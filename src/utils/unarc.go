package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var dest string = "./__files"

func Unarc(zipLink string) {
	arc, err := zip.OpenReader(zipLink)
	if err != nil {
		fmt.Println("Can't open archive")
	}
	defer arc.Close()
	for _, f := range arc.File {
		filePath := filepath.Join(dest, f.Name)
		isDir := f.FileInfo().IsDir()
		if isDir {
			err := os.MkdirAll(filePath, os.ModePerm)
			if err != nil {
				fmt.Printf("Error when create dir %v", f.Name)
			}
			continue
		}
		fmt.Printf("%v %v \n", f.Name, f.FileInfo().IsDir())
		err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
		if err != nil {
			return
		}

		srcFile, err := f.Open()
		if err != nil {
			return
		}
		defer srcFile.Close()

		destFile, err := os.Create(filePath)
		if err != nil {
			return
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return
		}
	}
}
