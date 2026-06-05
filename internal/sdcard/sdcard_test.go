package sdcard

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalLsblk(t *testing.T) {
	jsonData := `{
   "blockdevices": [
      {
         "name": "sda",
         "size": 512000000000,
         "mountpoint": null,
         "model": "Samsung SSD",
         "children": [
            {
               "name": "sda1",
               "size": 512000000,
               "mountpoint": "/boot",
               "model": null
            }
         ]
      }
   ]
}`

	var lsblk LsblkOutput
	err := json.Unmarshal([]byte(jsonData), &lsblk)
	if err != nil {
		t.Fatalf("Failed to unmarshal lsblk output: %v", err)
	}

	if len(lsblk.BlockDevices) != 1 {
		t.Fatalf("Expected 1 block device, got %d", len(lsblk.BlockDevices))
	}

	dev := lsblk.BlockDevices[0]
	if dev.Name != "sda" {
		t.Errorf("Expected name sda, got %s", dev.Name)
	}

	size, _ := dev.Size.Int64()
	if size != 512000000000 {
		t.Errorf("Expected size 512000000000, got %d", size)
	}

	if dev.MountPoint != nil {
		t.Errorf("Expected mountpoint to be nil, got %s", *dev.MountPoint)
	}

	if dev.Model == nil || *dev.Model != "Samsung SSD" {
		t.Errorf("Expected model Samsung SSD, got %v", dev.Model)
	}

	if len(dev.Children) != 1 {
		t.Fatalf("Expected 1 child, got %d", len(dev.Children))
	}

	child := dev.Children[0]
	if child.MountPoint == nil || *child.MountPoint != "/boot" {
		t.Errorf("Expected child mountpoint /boot, got %v", child.MountPoint)
	}
}
