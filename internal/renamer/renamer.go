package renamer

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
