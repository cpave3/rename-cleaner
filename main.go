package main

import (
	"fmt"
)

func main() {

	renamer := NewFileRenamer()
	renamer.parseFlags()

	err := renamer.scanFiles()
	if err != nil {
		fmt.Println("Error walking:", err)
	}

	if renamer.hasChanges() {

		if confirmed := renamer.confirmChanges(); confirmed {
			if err := renamer.executeRenames(); err != nil {
				fmt.Println("Error executing renames:", err)
				return
			}
		} else {
			fmt.Println("No changes made.")
		}
	} else {
		fmt.Println("No invalid files found. Goodbye!")
	}
}
