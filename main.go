package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/szogoon/piccp/internal/config"
	"github.com/szogoon/piccp/internal/copier"
	"github.com/szogoon/piccp/internal/grouping"
	"github.com/szogoon/piccp/internal/logging"
	"github.com/szogoon/piccp/internal/metadata"
	"github.com/szogoon/piccp/internal/progress"
	"github.com/szogoon/piccp/internal/sdcard"
)

var Version = "dev"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "piccp is a CLI utility for importing photos and videos from SD cards on Linux.\n\nUsage:\n")
		flag.PrintDefaults()
	}

	configPath := flag.String("config", "", "path to config file")
	dryRun := flag.Bool("dry-run", false, "dry run")
	versionFlag := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("piccp version %s\n", Version)
		os.Exit(0)
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		logging.Fatal("Failed to load config: %v", err)
	}

	if *dryRun {
		cfg.UI.DryRun = true
	}
	logging.Verbose = cfg.UI.Verbose

	if err := run(cfg); err != nil {
		// Fatal already exits with code 1, but we might want specific exit codes as per requirement
		logging.Error("%v", err)
		os.Exit(1)
	}
}

func run(cfg *config.Config) error {
	card, err := detectSDCard(cfg)
	if err != nil {
		return err
	}

	files, err := scanFiles(card, cfg)
	if err != nil {
		return err
	}

	logging.Info("Found %d files. Grouping into trips...", len(files))
	trips := grouping.GroupByTrip(files, cfg.Import.GroupingGapDays)
	logging.Info("Grouped into %d trips.", len(trips))

	if err := performImport(cfg, trips, files); err != nil {
		return err
	}

	logging.Info("Import completed successfully.")

	if cfg.SDCard.AutoUnmount && !cfg.UI.DryRun {
		logging.Info("Unmounting SD card...")
		if err := card.Unmount(); err != nil {
			logging.Error("Failed to unmount SD card: %v", err)
		}
	}

	return nil
}

func detectSDCard(cfg *config.Config) (*sdcard.SDCard, error) {
	logging.Info("Detecting SD cards...")
	cards, err := sdcard.ListSDCards(&cfg.SDCard)
	if err != nil {
		return nil, fmt.Errorf("error detecting SD cards: %w", err)
	}

	if len(cards) == 0 {
		return nil, fmt.Errorf("no SD cards found matching criteria")
	}

	// For MVP, just pick the first one
	card := &cards[0]
	logging.Info("Found SD card: %s mounted at %s", card.DevicePath, card.MountPoint)
	return card, nil
}

func scanFiles(card *sdcard.SDCard, cfg *config.Config) ([]metadata.FileInfo, error) {
	logging.Info("Scanning for files...")
	files, err := metadata.ScanFiles(card.MountPoint, cfg.Import.AllowedExtensions)
	if err != nil {
		return nil, fmt.Errorf("error scanning files: %w", err)
	}

	if len(files) == 0 {
		logging.Info("No supported files found.")
		os.Exit(0)
	}
	return files, nil
}

func performImport(cfg *config.Config, trips []grouping.Trip, files []metadata.FileInfo) error {
	var totalSize int64
	for _, f := range files {
		totalSize += f.Size
	}

	c := copier.NewCopier(cfg)

	var pb *progress.ProgressBar
	if cfg.UI.ProgressBar && !cfg.UI.DryRun {
		pb = progress.NewProgressBar(totalSize, "Importing")
	}

	err := c.ImportTrips(trips, func(n int64) {
		if pb != nil {
			pb.Add64(n)
		}
	})

	if pb != nil {
		pb.Finish()
		fmt.Println()
	}

	if err != nil {
		return fmt.Errorf("import failed: %w", err)
	}

	return nil
}
