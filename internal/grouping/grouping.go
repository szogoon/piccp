package grouping

import (
	"sort"
	"time"

	"github.com/szogoon/piccp/internal/metadata"
)

type Trip struct {
	Date  time.Time
	Files []metadata.FileInfo
}

func GroupByTrip(files []metadata.FileInfo, gapDays int) []Trip {
	if len(files) == 0 {
		return nil
	}

	// Sort files by timestamp
	sort.Slice(files, func(i, j int) bool {
		return files[i].Timestamp.Before(files[j].Timestamp)
	})

	var trips []Trip
	if len(files) == 0 {
		return trips
	}

	gap := time.Duration(gapDays) * 24 * time.Hour
	currentTrip := Trip{
		Date:  files[0].Timestamp,
		Files: []metadata.FileInfo{files[0]},
	}

	for i := 1; i < len(files); i++ {
		if files[i].Timestamp.Sub(files[i-1].Timestamp) > gap {
			trips = append(trips, currentTrip)
			currentTrip = Trip{
				Date:  files[i].Timestamp,
				Files: []metadata.FileInfo{files[i]},
			}
		} else {
			currentTrip.Files = append(currentTrip.Files, files[i])
		}
	}
	trips = append(trips, currentTrip)

	return trips
}
