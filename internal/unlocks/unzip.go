package unlocks

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func unzip(archivePath, targetDir string) error {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, f := range reader.File {
		parts := strings.SplitN(f.Name, "/", 2)
		if len(parts) < 2 {
			return fmt.Errorf("unsupported zip directory layout: %s", f.Name)
		}

		fpath := filepath.Join(targetDir, filepath.FromSlash(parts[1]))
		if !strings.HasPrefix(fpath, filepath.Clean(targetDir)+string(os.PathSeparator)) {
			continue
		}

		info := f.FileInfo()

		if info.IsDir() {
			err := os.MkdirAll(fpath, os.ModeDir|0700)
			if err != nil {
				return err
			}
		} else {
			err := os.MkdirAll(filepath.Dir(fpath), os.ModeDir|0700)
			if err != nil {
				return err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				return err
			}

			r, err := f.Open()
			if err != nil {
				outFile.Close()
				return err
			}

			_, err = io.Copy(outFile, r)
			if err != nil {
				outFile.Close()
				r.Close()
				return err
			}

			outFile.Close()
			r.Close()
		}
	}

	return nil
}
