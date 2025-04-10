package loader

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"vsixdm/config"
)

type ExtensionData []string

type PackageJSON struct {
	ExtensionPack []string `json:"extensionPack"`
}

func GetExtensionById(extensionId string) (*[]string, error) {
	cachedPaths, err := getCachedVSIXByExtId(extensionId, true)
	if err != nil {
		fmt.Printf("No extension in cache %q. Downloading from store.", extensionId)
	}
	if len(*cachedPaths) > 0 {
		return cachedPaths, nil
	}

	vsCodeMarketError := getExtensionFromCodeMarketPlace(extensionId)
	if vsCodeMarketError != nil {
		return &[]string{}, vsCodeMarketError
	}
	extensionPaths, err := getExtensionsListFolders(extensionId)
	if err != nil {
		return extensionPaths, err
	}

	var vsixList []string
	for _, extensionPath := range *extensionPaths {
		vsixPath, err := packExtensionToVSIX(extensionPath)
		if err != nil {
			return &vsixList, err
		}
		vsixList = append(vsixList, vsixPath)
	}
	return &vsixList, nil
}

func getExtensionsListFolders(extensionId string) (*[]string, error) {
	files, err := os.ReadDir(config.Paths.Extensions)
	if err != nil {
		return &[]string{}, err
	}

	var extensionPath string
	for _, file := range files {
		isExtensionFolder := file.IsDir() && strings.Contains(file.Name(), extensionId)
		if isExtensionFolder {
			extensionPath = filepath.Join(config.Paths.Extensions, file.Name())
		}
	}

	if len(extensionPath) == 0 {
		return &[]string{}, fmt.Errorf("extension folder not found %q", extensionId)
	}

	packageJsonPath := filepath.Join(extensionPath, "package.json")
	file, err := os.ReadFile(packageJsonPath)
	if err != nil {
		return &[]string{}, err
	}

	var pkgJSON PackageJSON
	parseError := json.Unmarshal([]byte(file), &pkgJSON)
	if parseError != nil {
		return &[]string{}, parseError
	}

	var extensionPackPaths []string
	extensionPackPaths = append(extensionPackPaths, extensionPath)
	for _, extensionId := range pkgJSON.ExtensionPack {
		for _, file := range files {
			isExtensionFolder := file.IsDir() && strings.Contains(file.Name(), extensionId)
			if isExtensionFolder {
				extensionPackPaths = append(extensionPackPaths, filepath.Join(config.Paths.Extensions, file.Name()))
			}

		}
	}

	return &extensionPackPaths, nil
}

func getExtensionFromCodeMarketPlace(extensionId string) error {
	cmd := exec.Command(config.Paths.VsCodeExec, "--install-extension", extensionId)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("can't call 'code --install extension %s'. Error: %v", extensionId, err)
	}

	cmd.Wait()
	lines := strings.Split(strings.Trim(string(output), "\n"), "\n")
	result := lines[len(lines)-1]

	if strings.Contains(result, "Failed") {
		return fmt.Errorf("can't get extension from marketplace %q. Err %v", extensionId, result)
	}
	return nil
}

func packExtensionToVSIX(extensionPath string) (string, error) {
	extensionFolder, extDirName := filepath.Split(extensionPath)
	tempDirName := fmt.Sprintf("temp-%s", strconv.FormatInt(time.Now().UnixNano(), 10))
	tempPath := filepath.Join(extensionFolder, tempDirName)
	err := os.Rename(extensionPath, tempPath)

	if err != nil {
		return "", err
	}

	errCreateFolder := os.Mkdir(extensionPath, os.ModeDir)
	if errCreateFolder != nil {
		return "", errCreateFolder
	}
	mvTempToExtFolderError := os.Rename(tempPath, filepath.Join(extensionPath, "extension"))
	if mvTempToExtFolderError != nil {
		return "", mvTempToExtFolderError
	}
	mkOutDirErr := os.MkdirAll(config.Paths.VSIXOutputDir, fs.ModeDir)
	if mkOutDirErr != nil {
		return "", err
	}
	vsixOutputPath := filepath.Join(config.Paths.VSIXOutputDir, fmt.Sprintf("%s.vsix", extDirName))

	zipFile, err := os.Create(vsixOutputPath)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	writer := zip.NewWriter(zipFile)
	defer writer.Close()

	errorAddToArchive := filepath.Walk(extensionPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == extensionPath {
			return nil
		}

		relPath := strings.TrimPrefix(path, extensionPath)
		relPath = strings.TrimPrefix(relPath, string(os.PathSeparator))
		if info.IsDir() {
			_, err := writer.Create(relPath + "/")
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("can't open file %v", err)
			return err
		}
		defer file.Close()
		fw, err := writer.Create(relPath)
		if err != nil {
			fmt.Printf("Writer error %v", err)
			return err
		}
		_, err = io.Copy(fw, file)
		if err != nil {
			fmt.Print(err)
		}
		return err
	})

	if errorAddToArchive != nil {
		fmt.Printf("Error when add to archive")
		return "", errorAddToArchive
	}
	rmError := os.RemoveAll(extensionPath)
	if rmError != nil {
		return "", rmError
	}
	return vsixOutputPath, nil

}

func getChildextensionsFromVSIX(VSIXPath string) (list *[]string, err error) {
	var childExtensions []string

	target := filepath.Join("extension", "package.json")
	r, err := zip.OpenReader(VSIXPath)
	if err != nil {
		return &childExtensions, err
	}

	defer r.Close()

	for _, file := range r.File {
		if file.Name == target {
			fileContent, err := file.Open()
			if err != nil {
				return &childExtensions, err
			}
			defer fileContent.Close()
			data, err := io.ReadAll(fileContent)
			if err != nil {
				return &childExtensions, err
			}
			var packageJson PackageJSON
			jsonParsingError := json.Unmarshal(data, &packageJson)
			if jsonParsingError != nil {
				return &childExtensions, jsonParsingError
			}
			childExtensions = packageJson.ExtensionPack
		}
	}
	return &childExtensions, nil
}

func getCachedVSIXByExtId(extensionId string, checkForNestedExtensions bool) (extPaths *[]string, err error) {
	files, err := os.ReadDir(config.Paths.VSIXOutputDir)
	if err != nil {
		if os.IsNotExist(err) {
			return &[]string{}, nil
		}
		return &[]string{}, err
	}

	var vsixPaths []string
	for _, file := range files {
		if strings.Contains(file.Name(), extensionId) {
			vsixPaths = append(vsixPaths, filepath.Join(config.Paths.VSIXOutputDir, file.Name()))
		}
	}
	if vsixPaths == nil {
		return &vsixPaths, fmt.Errorf("*.vsix extension for %q not found", extensionId)
	}
	if checkForNestedExtensions {
		childExts, err := getChildextensionsFromVSIX(vsixPaths[0])
		if err != nil {
			return &vsixPaths, err
		}
		for _, ext := range *childExts {
			extPath, err := getCachedVSIXByExtId(ext, false)
			if err != nil {
				return &vsixPaths, err
			}
			vsixPaths = append(vsixPaths, *extPath...)
		}
	}
	return &vsixPaths, nil
}
