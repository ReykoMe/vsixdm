package utils

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func copyFile(src string, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Printf("Can't open source file %q: %v\n", src, err)
		return err
	}
	defer srcFile.Close()

	targetFile, err := os.Create(dest)
	if err != nil {
		fmt.Printf("Can't create target file %q: %v\n", dest, err)
		return err
	}
	defer targetFile.Close()

	_, err = io.Copy(targetFile, srcFile)
	if err != nil {
		fmt.Printf("Can't copy content from %q to %q: %v\n", src, dest, err)
		return err
	}
	return nil
}

func copyDir(src string, dest string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Can't access path %q: %v\n", path, err)
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			fmt.Printf("Can't get relative path for %q: %v\n", path, err)
			return err
		}

		destPath := filepath.Join(dest, relPath)
		fmt.Println(relPath)
		if d.IsDir() {
			err := os.MkdirAll(destPath, d.Type().Perm())
			if err != nil {
				fmt.Printf("Can't create folder %q: %v\n", destPath, err)
				return err
			}
			return nil
		}
		return copyFile(path, destPath)
	})
}

func Copy(src string, dest string) error {
	entry, err := os.Stat(src)
	if err != nil {
		return err
	}
	if entry.IsDir() {
		return copyDir(src, dest)
	}
	return copyFile(src, dest)
}
