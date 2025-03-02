package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const validPattern = `^[a-zA-Z0-9_\-\.]+$`
const invalidPattern = `[^a-zA-Z0-9_\-\.]`

type FileRenamer struct {
	root   string
	dryRun bool
	files  []RenameItem
	dirs   []RenameItem
}

type RenameItem struct {
	oldPath string
	newPath string
}

func (fr *FileRenamer) scanFiles() error {
	return filepath.Walk(fr.root, fr.processPath)
}

func (fr *FileRenamer) processPath(path string, info os.FileInfo, err error) error {
	if err != nil {
		return fmt.Errorf("error accessing %s: %w", path, err)
	}

	filename := info.Name()

	if !isValidName(filename) {
		sanitizedName := sanitizeName(filename)
		newPath := filepath.Join(filepath.Dir(path), sanitizedName)
		fmt.Println("[Invalid] ", filename, " -> ", sanitizedName)

		// we need to do the dirs last, after the files have been renamed
		if info.IsDir() {
			fr.dirs = append(fr.dirs, struct {
				oldPath string
				newPath string
			}{path, newPath})
		} else {
			fr.files = append(fr.files, struct {
				oldPath string
				newPath string
			}{path, newPath})
		}
	}
	return nil
}

func (fr *FileRenamer) hasChanges() bool {
	return len(fr.files) > 0 || len(fr.dirs) > 0
}

func (fr *FileRenamer) confirmChanges() bool {
	fmt.Println()
	if fr.dryRun {
		fmt.Println("DRY RUN: No changes will be made.")
	}

	fmt.Println("Do you want to rename these files? (y/N)")
	var response string
	_, err := fmt.Scanln(&response)
	fmt.Println()
	if err != nil {
		fmt.Println("Error reading input:", err)
		return false
	}

	return strings.ToLower(response) == "y"
}

func (fr *FileRenamer) renameItem(item RenameItem) error {
	if fr.dryRun {
		fmt.Println("[DRY-RUN] Would rename:", item.oldPath, item.newPath)
	} else {
		err := os.Rename(item.oldPath, item.newPath)
		if err != nil {
			fmt.Println("Error renaming:", item.oldPath, item.newPath, err)
			return err
		} else {
			fmt.Println("Renamed:", item.oldPath, item.newPath)
		}
	}
	return nil
}

func (fr *FileRenamer) executeRenames() error {
	for _, item := range fr.files {
		if err := fr.renameItem(item); err != nil {
			return err
		}
	}

	// Rename the directories in reverse, to avoid conflicts
	for i := len(fr.dirs) - 1; i >= 0; i-- {
		item := fr.dirs[i]
		if err := fr.renameItem(item); err != nil {
			return err
		}
	}

	fr.dirs = nil
	fr.files = nil
	return nil
}

func (fr *FileRenamer) parseFlags() {
	root := flag.String("root", ".", "Root directory to scan")
	dryRun := flag.Bool("dry-run", false, "Preview changes without making them")
	flag.Parse()

	if len(flag.Args()) > 0 {
		*root = flag.Args()[0]
	}

	if _, err := os.Stat(*root); os.IsNotExist(err) {
		fmt.Println("Error: Root directory does not exist:", *root)
		os.Exit(1)
	}

	fr.root = *root
	fr.dryRun = *dryRun
}

func NewFileRenamer() FileRenamer {
	return FileRenamer{
		".",
		false,
		[]RenameItem{},
		[]RenameItem{},
	}
}

func isValidName(name string) bool {
	// Only allow alphanumeric characters, hyphens, and periods.
	validPattern := regexp.MustCompile(validPattern)
	return validPattern.MatchString(name)
}

func sanitizeName(name string) string {
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "[", "-")
	name = strings.ReplaceAll(name, "]", "-")

	invalidChars := regexp.MustCompile(invalidPattern)
	name = invalidChars.ReplaceAllString(name, "")
	name = regexp.MustCompile(`-{2,}`).ReplaceAllString(name, "-")
	name = strings.Trim(name, "-")

	// If there is a non alphanumeric character before the final period, remove it.
	// i.e. file[name].txt -> file-name-.txt -> file-name.txt
	if strings.Contains(name, ".") {
		lastPeriodIndex := strings.LastIndex(name, ".")
		priorCharIndex := lastPeriodIndex - 1
		priorChar := name[priorCharIndex]
		// if prior chara is not alphanumeric, remove it
		if regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(string(priorChar)) {
			// remove the character at the prior index
			name = name[:priorCharIndex] + name[lastPeriodIndex:]
		}

	}

	return name
}
