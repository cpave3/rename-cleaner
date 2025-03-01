# Rename Files Program

A simple Go program to scan a directory for files and folders that don't match a desired naming pattern, list them, and rename them upon user confirmation.
Features

Scans directories for files and folders that don't match a predefined naming pattern.
Lists the items that would be renamed.
Renames only those files and folders upon user confirmation.
Option to run in "dry-run" mode to preview changes without renaming.

## Installation

Clone this repository:
git clone https://github.com/cpave3/rename-files.git

Install Go dependencies:
`go mod tidy`

Build the binary:
`go build -o rename-files`

Run the program:
`./rename-files`

## Usage

./rename-files --dry-run /path/to/your/directory

- NOTE: currently the dry run flag should come before the directory path.
- Replace /path/to/your/directory with the directory to scan.
- Use the --dry-run flag to preview changes without renaming.

## License

MIT License.