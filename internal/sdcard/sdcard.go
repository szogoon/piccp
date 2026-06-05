package sdcard

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/szogoon/piccp/internal/config"
)

type BlockDevice struct {
	Name       string        `json:"name"`
	Size       string        `json:"size"`
	MountPoint string        `json:"mountpoint"`
	Model      string        `json:"model"`
	Children   []BlockDevice `json:"children"`
}

type LsblkOutput struct {
	BlockDevices []BlockDevice `json:"blockdevices"`
}

type SDCard struct {
	DevicePath string
	MountPoint string
	SizeGB     float64
}

func ListSDCards(cfg *config.SDCardConfig) ([]SDCard, error) {
	cmd := exec.Command("lsblk", "-J", "-b", "-o", "NAME,SIZE,MOUNTPOINT,MODEL")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run lsblk: %w", err)
	}

	var lsblk LsblkOutput
	if err := json.Unmarshal(output, &lsblk); err != nil {
		return nil, fmt.Errorf("failed to parse lsblk output: %w", err)
	}

	var cards []SDCard
	for _, dev := range lsblk.BlockDevices {
		// Look for removable devices or devices with specific models if needed
		// For now, let's filter by size and check children for mountpoints
		processDevice(dev, cfg, &cards)
	}

	return cards, nil
}

func processDevice(dev BlockDevice, cfg *config.SDCardConfig, cards *[]SDCard) {
	sizeBytes, _ := strconv.ParseFloat(dev.Size, 64)
	sizeGB := sizeBytes / (1024 * 1024 * 1024)

	if sizeGB >= cfg.MinSizeGB && sizeGB <= cfg.MaxSizeGB {
		// Check if it's mounted or has mounted partitions
		mountPoint := dev.MountPoint
		if mountPoint == "" && len(dev.Children) > 0 {
			for _, child := range dev.Children {
				if child.MountPoint != "" {
					mountPoint = child.MountPoint
					break
				}
			}
		}

		if mountPoint != "" {
			*cards = append(*cards, SDCard{
				DevicePath: "/dev/" + dev.Name,
				MountPoint: mountPoint,
				SizeGB:     sizeGB,
			})
		}
	}

	for _, child := range dev.Children {
		processDevice(child, cfg, cards)
	}
}

// Mount and Unmount placeholders for now
func (s *SDCard) Unmount() error {
	cmd := exec.Command("umount", s.MountPoint)
	return cmd.Run()
}
