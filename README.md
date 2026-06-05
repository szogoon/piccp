# piccp

`piccp` is a CLI utility for importing photos and videos from SD cards on Linux.

## Features

- Detects SD cards using `lsblk`.
- Groups files into "trips" based on modification date continuity.
- Copies files to structured directories (`YYYY-MM-DD`).
- Global progress bar for imports.
- Dry-run mode for safety.

## Installation

```bash
go install github.com/szogoon/piccp@latest
```

## Usage

```bash
piccp --config ~/.config/piccp/config.toml
```

## Configuration (TOML)

Default location: `~/.config/piccp/config.toml`

```toml
[target]
output_dir = "/home/user/Pictures/Imports"

[sdcard]
min_size_gb = 4
max_size_gb = 1024
auto_unmount = true

[import]
grouping_gap_days = 1
allowed_extensions = ["jpg", "jpeg", "png", "mp4", "mov"]
preserve_directory_structure = false
overwrite_existing = false

[ui]
progress_bar = true
verbose = true
dry_run = false
```
