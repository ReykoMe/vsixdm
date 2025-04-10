package config

import (
	"os"
	"path/filepath"
)

type AppPaths struct {
	Root          string
	Extensions    string
	VsCodeExec    string
	VSIXOutputDir string
	Translations  string
}

var Paths *AppPaths = func() *AppPaths {
	stage := os.Getenv("VSIXDM_STAGE")

	__rootDir, err := os.Executable()
	__rootDir = filepath.Dir(__rootDir)

	if stage == "DEV" {
		__rootDir, err = os.Getwd()
	}

	if err != nil {
		return &AppPaths{}
	}

	__vscodePortDir := filepath.Join(__rootDir, "code")

	if stage == "DEV" {
		__vscodePortDir = filepath.Join(__rootDir, "..", "vscode", "win")
	}

	appPaths := &AppPaths{
		Root:          __rootDir,
		Extensions:    filepath.Join(__vscodePortDir, "data", "extensions"),
		VsCodeExec:    filepath.Join(__vscodePortDir, "code.cmd"),
		VSIXOutputDir: filepath.Join(__rootDir, "output"),
		Translations:  filepath.Join(__rootDir, "translations"),
	}
	if stage == "DEV" {
		appPaths.VSIXOutputDir = filepath.Join(__rootDir, "..", "output")
	}
	return appPaths
}()
