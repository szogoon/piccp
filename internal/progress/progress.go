package progress

import (
	"github.com/schollz/progressbar/v3"
)

type ProgressBar = progressbar.ProgressBar

func NewProgressBar(total int64, description string) *ProgressBar {
	return progressbar.DefaultBytes(
		total,
		description,
	)
}
