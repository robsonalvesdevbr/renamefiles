package main

import (
	"fmt"
	"os"

	"github.com/robsonalvesdevbr/renamefiles/internal/renamer"
	"github.com/spf13/cobra"
)

func main() {
	var useUnderscore bool
	var removeUnderscore bool
	var separator string
	var oldSeparator string
	var newSeparator string
	var toTitleCase bool
	var includeTimestamp bool
	var dryRun bool

	var rootCmd = &cobra.Command{
		Use:   "renamefiles",
		Short: "RenameFiles is a tool to rename files in a directory",
		Run: func(cmd *cobra.Command, args []string) {
			if useUnderscore && removeUnderscore {
				fmt.Println("Erro: Não é possível usar 'underscore' e 'remove-underscore' ao mesmo tempo.")
				os.Exit(1)
			}
			dir, err := os.Getwd()
			if err != nil {
				fmt.Printf("Error getting current directory: %v\n", err)
				os.Exit(1)
			}
			fileRenamer := renamer.NewFileRenamer(separator, useUnderscore, removeUnderscore, oldSeparator, newSeparator, toTitleCase, includeTimestamp)
			if err := fileRenamer.RenameFiles(dir, dryRun); err != nil {
				fmt.Printf("Error renaming files: %v\n", err)
				os.Exit(1)
			}
		},
	}

	rootCmd.Flags().BoolVar(&useUnderscore, "underscore", false, "Replace spaces with underscores in file names")
	rootCmd.Flags().BoolVar(&removeUnderscore, "remove-underscore", false, "Replace underscores with spaces in file names")
	rootCmd.Flags().StringVar(&separator, "separator", "", "Character to use as a separator (e.g., _ or -)")
	rootCmd.Flags().StringVar(&oldSeparator, "old-separator", "", "Character to be replaced in file names")
	rootCmd.Flags().StringVar(&newSeparator, "new-separator", "", "Character to replace the old separator in file names")
	rootCmd.Flags().BoolVar(&toTitleCase, "title-case", false, "Convert file names to Title Case (capitalize first letter of each word)")
	rootCmd.Flags().BoolVar(&includeTimestamp, "include-timestamp", false, "Include the file's creation timestamp in the file name")
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Perform a dry run without renaming files")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
