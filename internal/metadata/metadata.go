package metadata

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileInfo struct {
	Path      string
	Timestamp time.Time
	Size      int64
}

func ScanFiles(root string, allowedExtensions []string) ([]FileInfo, error) {
	var files []FileInfo
	extMap := make(map[string]bool)
	for _, ext := range allowedExtensions {
		extMap[strings.ToLower(ext)] = true
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if len(ext) > 0 {
			ext = ext[1:] // remove dot
		}

		if extMap[ext] {
			files = append(files, FileInfo{
				Path:      path,
				Timestamp: info.ModTime(),
				Size:      info.Size(),
			})
		}
		return nil
	})

	return files, err
}
