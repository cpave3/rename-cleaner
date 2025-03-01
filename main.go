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

func main() {

	var renameList []struct {
		oldPath string
		newPath string
	}

	var dirRenameList []struct {
		oldPath string
		newPath string
	}

	root := flag.String("root", ".", "Root directory to scan")
	dryRun := flag.Bool("dry-run", false, "Preview changes without making them")
	flag.Parse()

	if len(flag.Args()) > 0 {
		*root = flag.Args()[0]
	}

	if _, err := os.Stat(*root); os.IsNotExist(err) {
		fmt.Println("Error: Root directory does not exist:", *root)
		return
	}

	// Recursively walk the directory tree, from where we are executed.
	err := filepath.Walk(*root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error accessing:", path, err)
			return nil
		}
		filename := info.Name()
		if !isValidName(filename) {
			sanitizedName := sanitizeName(filename)
			newPath := filepath.Join(filepath.Dir(path), sanitizedName)
			fmt.Println("[Invalid] ", filename, " -> ", sanitizedName)

			// we need to do the dirs last, after the files have been renamed
			if info.IsDir() {
				dirRenameList = append(dirRenameList, struct {
					oldPath string
					newPath string
				}{path, newPath})
			} else {
				renameList = append(renameList, struct {
					oldPath string
					newPath string
				}{path, newPath})
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking:", err)
	}

	if len(renameList) > 0 || len(dirRenameList) > 0 {

		fmt.Println()
		if *dryRun {
			fmt.Println("DRY RUN: No changes will be made.")
		}

		fmt.Println("Do you want to rename these files? (y/N)")
		var response string
		_, err := fmt.Scanln(&response)
		fmt.Println()
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}

		if strings.ToLower(response) == "y" {

			for _, item := range renameList {
				if *dryRun {
					fmt.Println("[DRY-RUN] Would rename:", item.oldPath, item.newPath)
				} else {
					err := os.Rename(item.oldPath, item.newPath)
					if err != nil {
						fmt.Println("Error renaming:", item.oldPath, item.newPath, err)
						return
					} else {
						fmt.Println("Renamed:", item.oldPath, item.newPath)
					}
				}
			}

			// Rename the directories in reverse, to avoid conflicts
			for i := len(dirRenameList) - 1; i >= 0; i-- {
				item := dirRenameList[i]
				if *dryRun {
					fmt.Println("[DRY-RUN] Would rename:", item.oldPath, item.newPath)
				} else {
					err := os.Rename(item.oldPath, item.newPath)
					if err != nil {
						fmt.Println("Error renaming:", item.oldPath, item.newPath, err)
						return
					} else {
						fmt.Println("Renamed:", item.oldPath, item.newPath)
					}
				}

			}
		} else {
			fmt.Println("No changes made. Goodbye!")
		}
	} else {
		fmt.Println("No invalid files found. Goodbye!")
	}
}

func isValidName(name string) bool {
	// Only allow alphanumeric characters, hyphens, and periods.
	validPattern := regexp.MustCompile(validPattern)
	return validPattern.MatchString(name)
}

func sanitizeName(name string) string {
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "[", "")
	name = strings.ReplaceAll(name, "]", "")

	invalidChars := regexp.MustCompile(invalidPattern)
	name = invalidChars.ReplaceAllString(name, "")
	name = regexp.MustCompile(`-{2,}`).ReplaceAllString(name, "-")
	name = strings.Trim(name, "-")
	return name
}
