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
func sanitizeFileName(name string, useUnderscore bool, removeUnderscore bool) string {
	// Normalize Unicode to remove stylization.
	sanitized := normalizeUnicode(name)

	if useUnderscore {
		// Replace spaces with underscores.
		sanitized = strings.ReplaceAll(sanitized, " ", "_")
	}

	if removeUnderscore {
		// Replace underscores with spaces.
		sanitized = strings.ReplaceAll(sanitized, "_", " ")
	}

	// Trim extra spaces.
	sanitized = strings.TrimSpace(sanitized)

	// Replace invalid characters with underscore using regex (excluding special characters).
	invalidChars := regexp.MustCompile(`[^ -~ -ÿ\p{L}\p{N}_.-]`)
	sanitized = invalidChars.ReplaceAllString(sanitized, "_")

	return sanitized
}

// renameFiles iterates through the files in the directory and renames them.
func renameFiles(dir string, useUnderscore bool, removeUnderscore bool) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			oldName := d.Name()
			newName := sanitizeFileName(oldName, useUnderscore, removeUnderscore)

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
	flag.Parse()

	// Get the current working directory.
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Current directory: %s\n", dir)

	// Rename files in the current directory.
	if err := renameFiles(dir, *useUnderscore, *removeUnderscore); err != nil {
		fmt.Printf("Error renaming files: %v\n", err)
		os.Exit(1)
	}
}
