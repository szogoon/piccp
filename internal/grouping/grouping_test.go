package grouping

import (
	"testing"
	"time"

	"github.com/szogoon/piccp/internal/metadata"
)

func TestGroupByTrip(t *testing.T) {
	now := time.Now()
	day := 24 * time.Hour

	files := []metadata.FileInfo{
		{Path: "f1.jpg", Timestamp: now, Size: 100},
		{Path: "f2.jpg", Timestamp: now.Add(1 * time.Hour), Size: 100},
		{Path: "f3.jpg", Timestamp: now.Add(2 * day), Size: 100},
		{Path: "f4.jpg", Timestamp: now.Add(2 * day).Add(1 * time.Hour), Size: 100},
	}

	trips := GroupByTrip(files, 1)

	if len(trips) != 2 {
		t.Errorf("Expected 2 trips, got %d", len(trips))
	}

	if len(trips[0].Files) != 2 {
		t.Errorf("Expected 2 files in first trip, got %d", len(trips[0].Files))
	}

	if len(trips[1].Files) != 2 {
		t.Errorf("Expected 2 files in second trip, got %d", len(trips[1].Files))
	}
}

func TestGroupByTrip_Empty(t *testing.T) {
	trips := GroupByTrip(nil, 1)
	if trips != nil {
		t.Error("Expected nil trips for nil input")
	}
}

func TestGroupByTrip_Sorting(t *testing.T) {
	now := time.Now()
	files := []metadata.FileInfo{
		{Path: "f2.jpg", Timestamp: now.Add(1 * time.Hour), Size: 100},
		{Path: "f1.jpg", Timestamp: now, Size: 100},
	}

	trips := GroupByTrip(files, 1)
	if len(trips) != 1 {
		t.Fatalf("Expected 1 trip, got %d", len(trips))
	}

	if trips[0].Files[0].Path != "f1.jpg" {
		t.Errorf("Expected f1.jpg to be first after sorting, got %s", trips[0].Files[0].Path)
	}
}
