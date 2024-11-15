package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

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

// sanitizeFileName removes spaces and non-visible characters while keeping special characters.
func sanitizeFileName(name string, separator string, useUnderscore bool, removeUnderscore bool, oldSeparator string, newSeparator string) string {
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

	return sanitized
}

// renameFiles iterates through the files in the directory and renames them.
func renameFiles(dir string, separator string, useUnderscore bool, removeUnderscore bool, oldSeparator string, newSeparator string) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			oldName := d.Name()
			newName := sanitizeFileName(oldName, separator, useUnderscore, removeUnderscore, oldSeparator, newSeparator)

			if oldName != newName {
				oldPath := filepath.Join(dir, oldName)
				newPath := filepath.Join(dir, newName)

				if err := os.Rename(oldPath, newPath); err != nil {
					return fmt.Errorf("failed to rename file %s to %s: %w", oldName, newName, err)
				}

				fmt.Printf("Renamed: %s -> %s\n", oldName, newName)
			}
		}
		return nil
	})
}

func main() {
	// Parse command-line flags.
	useUnderscore := flag.Bool("underscore", false, "Replace spaces with underscores in file names")
	removeUnderscore := flag.Bool("remove-underscore", false, "Replace underscores with spaces in file names")
	separator := flag.String("separator", "", "Character to use as a separator (e.g., _ or -)")
	oldSeparator := flag.String("old-separator", "", "Character to be replaced in file names")
	newSeparator := flag.String("new-separator", "", "Character to replace the old separator in file names")
	flag.Parse()

	// Get the current working directory.
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Current directory: %s\n", dir)

	// Rename files in the current directory.
	if err := renameFiles(dir, *separator, *useUnderscore, *removeUnderscore, *oldSeparator, *newSeparator); err != nil {
		fmt.Printf("Error renaming files: %v\n", err)
		os.Exit(1)
	}
}
