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

// FileRenamer interface defines the contract for renaming files.
type FileRenamer interface {
	RenameFiles(dir string, dryRun bool) error
}

// FileRenamerImpl implements the FileRenamer interface.
type FileRenamerImpl struct {
	separator        string
	useUnderscore    bool
	removeUnderscore bool
	oldSeparator     string
	newSeparator     string
	toTitleCase      bool
	includeTimestamp bool
}

// NewFileRenamer creates a new instance of FileRenamerImpl.
func NewFileRenamer(separator string, useUnderscore, removeUnderscore bool, oldSeparator, newSeparator string, toTitleCase, includeTimestamp bool) FileRenamer {
	return &FileRenamerImpl{
		separator:        separator,
		useUnderscore:    useUnderscore,
		removeUnderscore: removeUnderscore,
		oldSeparator:     oldSeparator,
		newSeparator:     newSeparator,
		toTitleCase:      toTitleCase,
		includeTimestamp: includeTimestamp,
	}
}

// RenameFiles iterates through the files in the directory and renames them.
func (fr *FileRenamerImpl) RenameFiles(dir string, dryRun bool) error {
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
			newName := fr.sanitizeFileName(oldName, timestamp)
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

// sanitizeFileName removes spaces and non-visible characters while keeping special characters.
func (fr *FileRenamerImpl) sanitizeFileName(name, timestamp string) string {
	sanitized := normalizeUnicode(name)
	if fr.useUnderscore {
		sanitized = strings.ReplaceAll(sanitized, " ", "_")
	} else if fr.removeUnderscore {
		sanitized = strings.ReplaceAll(sanitized, "_", " ")
	} else if fr.separator != "" {
		sanitized = strings.ReplaceAll(sanitized, " ", fr.separator)
	}
	if fr.oldSeparator != "" && fr.newSeparator != "" {
		sanitized = strings.ReplaceAll(sanitized, fr.oldSeparator, fr.newSeparator)
	}
	sanitized = strings.TrimSpace(sanitized)
	invalidChars := regexp.MustCompile(`[^ -~\xa0-\xff\p{L}\p{N}_.-]`)
	if fr.separator != "" {
		sanitized = invalidChars.ReplaceAllString(sanitized, fr.separator)
	} else {
		sanitized = invalidChars.ReplaceAllString(sanitized, "")
	}
	if fr.toTitleCase {
		sanitized = toTitleCaseWithSeparator(sanitized, fr.separator)
	}

	// Verifica se o nome já contém o timestamp
	if fr.includeTimestamp {
		ext := filepath.Ext(sanitized)
		nameWithoutExt := strings.TrimSuffix(sanitized, ext)
		timestampPattern := regexp.MustCompile(`^\d{8}_\d{6}_`)
		if !timestampPattern.MatchString(nameWithoutExt) {
			sanitized = fmt.Sprintf("%s_%s%s", timestamp, nameWithoutExt, ext)
		}
	}

	return sanitized
}

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
			fileRenamer := NewFileRenamer(separator, useUnderscore, removeUnderscore, oldSeparator, newSeparator, toTitleCase, includeTimestamp)
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
