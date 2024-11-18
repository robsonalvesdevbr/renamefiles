package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/spf13/cobra"
	"golang.org/x/text/unicode/norm"
)

// normalizeUnicode removes diacritics and converts styled characters to their basic form.
func normalizeUnicode(input string) string {
	return strings.Map(func(r rune) rune {
		if unicode.Is(unicode.Mn, r) { // Remove diacritics
			return -1
		}
		return r
	}, norm.NFKD.String(input))
}

// toTitleCaseWithSeparator converts a string to Title Case while maintaining separators.
func toTitleCaseWithSeparator(input, separator string) string {
	if separator == "" {
		separator = " "
	}
	words := strings.Split(input, separator)
	for i, word := range words {
		if len(word) > 0 {
			if len(word) > 1 {
				words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
			} else {
				words[i] = strings.ToUpper(word)
			}
		}
	}
	return strings.Join(words, separator)
}

// sanitizeFileName removes spaces and non-visible characters while keeping special characters.
func sanitizeFileName(name string, separator string, useUnderscore bool, removeUnderscore bool, oldSeparator string, newSeparator string, toTitleCase bool, includeTimestamp bool, timestamp string) string {
	// Normalize Unicode to remove stylization.
	sanitized := normalizeUnicode(name)
	if useUnderscore {
		sanitized = strings.ReplaceAll(sanitized, " ", "_")
	} else if removeUnderscore {
		sanitized = strings.ReplaceAll(sanitized, "_", " ")
	} else if separator != "" {
		sanitized = strings.ReplaceAll(sanitized, " ", separator)
	}
	if oldSeparator != "" && newSeparator != "" {
		sanitized = strings.ReplaceAll(sanitized, oldSeparator, newSeparator)
	}
	// Trim extra spaces.
	sanitized = strings.TrimSpace(sanitized)
	// Replace invalid characters with the separator using regex (excluding special characters).
	invalidChars := regexp.MustCompile(`[^ -~ -ÿ\p{L}\p{N}_.-]`)
	if separator != "" {
		sanitized = invalidChars.ReplaceAllString(sanitized, separator)
	} else {
		sanitized = invalidChars.ReplaceAllString(sanitized, "")
	}
	if toTitleCase {
		sanitized = toTitleCaseWithSeparator(sanitized, separator)
	}
	if includeTimestamp {
		ext := filepath.Ext(sanitized)
		nameWithoutExt := strings.TrimSuffix(sanitized, ext)
		sanitized = fmt.Sprintf("%s_%s%s", timestamp, nameWithoutExt, ext)
	}
	return sanitized
}

// renameFiles iterates through the files in the directory and renames them.
func renameFiles(dir string, separator string, useUnderscore bool, removeUnderscore bool, oldSeparator string, newSeparator string, toTitleCase bool, includeTimestamp bool, dryRun bool) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			oldName := d.Name()
			info, err := d.Info()
			if err != nil {
				return err
			}
			timestamp := info.ModTime().Format("20060102_150405")
			newName := sanitizeFileName(oldName, separator, useUnderscore, removeUnderscore, oldSeparator, newSeparator, toTitleCase, includeTimestamp, timestamp)
			if oldName != newName {
				oldPath := path
				newPath := filepath.Join(filepath.Dir(path), newName)
				if dryRun {
					fmt.Printf("Dry run: %s -> %s\n", oldName, newName)
				} else {
					if err := os.Rename(oldPath, newPath); err != nil {
						return fmt.Errorf("failed to rename file %s to %s: %w", oldName, newName, err)
					}
					fmt.Printf("Renamed: %s -> %s\n", oldName, newName)
				}
			}
		}
		return nil
	})
}

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
			// Ensure conflicting flags are not used simultaneously.
			if useUnderscore && removeUnderscore {
				fmt.Println("Erro: Não é possível usar 'underscore' e 'remove-underscore' ao mesmo tempo.")
				os.Exit(1)
			}
			// Get the current working directory.
			dir, err := os.Getwd()
			if err != nil {
				fmt.Printf("Error getting current directory: %v\n", err)
				os.Exit(1)
			}
			// Rename files in the directory.
			if err := renameFiles(dir, separator, useUnderscore, removeUnderscore, oldSeparator, newSeparator, toTitleCase, includeTimestamp, dryRun); err != nil {
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
