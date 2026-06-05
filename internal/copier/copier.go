package copier

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/szogoon/piccp/internal/config"
	"github.com/szogoon/piccp/internal/grouping"
)

type Copier struct {
	Config *config.Config
}

func NewCopier(cfg *config.Config) *Copier {
	return &Copier{Config: cfg}
}

func (c *Copier) ImportTrips(trips []grouping.Trip, onProgress func(int64)) error {
	for _, trip := range trips {
		if err := c.importTrip(trip, onProgress); err != nil {
			return err
		}
	}
	return nil
}

func (c *Copier) importTrip(trip grouping.Trip, onProgress func(int64)) error {
	tripDirName := trip.Date.Format("2006-01-02")
	targetDir := filepath.Join(c.Config.Target.OutputDir, tripDirName)

	if !c.Config.UI.DryRun {
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", targetDir, err)
		}
	}

	for _, file := range trip.Files {
		destPath := filepath.Join(targetDir, filepath.Base(file.Path))
		if err := c.copyFile(file.Path, destPath, onProgress); err != nil {
			return err
		}
	}
	return nil
}

func (c *Copier) copyFile(src, dst string, onProgress func(int64)) error {
	if c.Config.UI.DryRun {
		fmt.Printf("[DRY-RUN] Copying %s to %s\n", src, dst)
		return nil
	}

	if c.shouldSkip(dst) {
		return nil
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	return c.copyData(source, destination, onProgress)
}

func (c *Copier) shouldSkip(dst string) bool {
	if !c.Config.Import.OverwriteExisting {
		if _, err := os.Stat(dst); err == nil {
			return true
		}
	}
	return false
}

func (c *Copier) copyData(src io.Reader, dst io.Writer, onProgress func(int64)) error {
	buf := make([]byte, 32*1024)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			if _, err := dst.Write(buf[:n]); err != nil {
				return err
			}
			if onProgress != nil {
				onProgress(int64(n))
			}
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}
