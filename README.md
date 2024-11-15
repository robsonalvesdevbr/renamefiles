# Golang File Renamer

This Go application provides a robust way to rename files in a directory. It includes several configurable options for modifying file names, such as replacing spaces with underscores, converting file names to Title Case, and customizing separators.

## Features

- Replace spaces with a specified separator.
- Remove underscores and replace them with spaces.
- Replace specific characters with others.
- Convert file names to Title Case (capitalize the first letter of each word).
- Fully customizable via command-line flags.

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd <repository-name>
   ```
2. Build the application:
   ```bash
   go build -o file-renamer
   ```

## Usage

Run the application from the directory containing the files to rename. Use the available flags to customize the behavior.

### Command-Line Flags

| Flag                  | Description                                                                  | Default Value |
| --------------------- | ---------------------------------------------------------------------------- | ------------- |
| `--underscore`        | Replace spaces with underscores in file names.                               | `false`       |
| `--remove-underscore` | Replace underscores with spaces in file names.                               | `false`       |
| `--separator`         | Specify a custom separator to replace spaces.                                | `""`          |
| `--old-separator`     | Character to be replaced in file names.                                      | `""`          |
| `--new-separator`     | Character to replace the old separator in file names.                        | `""`          |
| `--title-case`        | Convert file names to Title Case (capitalize the first letter of each word). | `false`       |

### Examples

1. **Replace spaces with underscores:**

   ```bash
   ./file-renamer --underscore
   ```

2. **Replace underscores with spaces:**

   ```bash
   ./file-renamer --remove-underscore
   ```

3. **Replace spaces with a custom separator (e.g., hyphen):**

   ```bash
   ./file-renamer --separator=-
   ```

4. **Replace specific characters:**

   ```bash
   ./file-renamer --old-separator=_ --new-separator=-
   ```

5. **Convert to Title Case while keeping separators:**

   ```bash
   ./file-renamer --title-case --separator=-
   ```

6. **Combine multiple options:**
   ```bash
   ./file-renamer --underscore --title-case
   ```

### Output Example

For a directory with the following files:

```
my_file.txt
another file.txt
```

Running:

```bash
./file-renamer --title-case --separator=_
```

Produces:

```
My_File.txt
Another_File.txt
```

## Error Handling

If an error occurs (e.g., permission issues, invalid file paths), the application will log the error and stop execution.

## Development

This application uses the Go standard library and the `golang.org/x/text/unicode/norm` package for Unicode normalization.

### Key Functions

- **`normalizeUnicode`**: Removes diacritics and converts styled characters to their basic form.
- **`toTitleCaseWithSeparator`**: Converts strings to Title Case while respecting separators.
- **`sanitizeFileName`**: Core function for applying all transformations to file names.
- **`renameFiles`**: Iterates through files in a directory and applies renaming based on the configured flags.

## License

This project is licensed under the MIT License. See the LICENSE file for details.
