package utee

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Zip compress file or directory into writer
// pathToZip: the path to zip, could be a file or a directory
// dest: the writer to write zip content
func Zip(pathToZip string, dest io.Writer) error {
	// create zip writer
	zipWriter := zip.NewWriter(dest)

	fn := func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		// relative path
		relPath := strings.TrimPrefix(filePath, filepath.Dir(pathToZip))
		zipFile, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		return err
	}

	if err := filepath.Walk(pathToZip, fn); err != nil {
		return err
	}

	err := zipWriter.Close()
	return err
}
