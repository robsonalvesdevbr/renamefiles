# RenameFiles

RenameFiles is a Go-based utility for renaming files in a directory, offering various options for normalizing file names, including converting spaces to underscores, removing diacritics, and changing letter cases., including converting spaces to underscores, removing diacritics, and changing letter cases.

## Features

- Normalize Unicode characters by removing diacritics.
- Replace spaces with custom separators (e.g., underscores or dashes).
- Convert file names to Title Case.
- Customize the way separators are used in file names.

## How It Works

RenameFiles scans the current directory and renames files based on the given options. You can apply different transformations to clean up or standardize the naming format of your files.. You can apply different transformations to clean up or standardize the naming format of your files.

## Command-Line Flags

The utility supports the following command-line flags to customize the behavior of the renaming process:

| Flag                                              | Description                                                                                                     |
| ------------------------------------------------- | --------------------------------------------------------------------------------------------------------------- |
| `-underscore`                                     | Replace spaces with underscores in file names.                                                                  |
| `-remove-underscore`                              | Replace underscores with spaces in file names.                                                                  |
| `-separator=<separator>`                          | Specify a character to use as a separator (e.g., `_` or `-`). This replaces spaces with the provided separator. |
| `-old-separator=<old>` and `-new-separator=<new>` | Replace the old separator character with a new one in file names.                                               |
| `-title-case`                                     | Convert file names to Title Case (capitalize the first letter of each word).                                    |

> **Note**: You cannot use `-underscore` and `-remove-underscore` simultaneously. The tool will exit with an error if this occurs.

## Usage

To use the tool, compile it and run it in the desired directory where your files are located. Here are a few examples:

### Replace Spaces with Underscores

```
$ go run main.go -underscore
```

This command replaces all spaces in file names with underscores (`_`).

### Replace Underscores with Spaces

```
$ go run main.go -remove-underscore
```

This command replaces all underscores in file names with spaces.

### Use a Custom Separator

```
$ go run main.go -separator=-
```

This command replaces all spaces in file names with hyphens (`-`).

### Change a Specific Separator

```
$ go run main.go -old-separator=_ -new-separator=-
```

This command replaces all underscores (`_`) in file names with hyphens (`-`).

### Convert File Names to Title Case

```
$ go run main.go -title-case
```

This command capitalizes the first letter of each word in the file names.

## Example Output

After running the tool, you will see output like the following:

```
Current directory: /path/to/your/directory
Renamed: old_file_name.txt -> Old_File_Name.txt
Renamed: another-file.txt -> Another_File.txt
```

## Error Handling

- The tool will skip renaming if the new file name already exists, to avoid overwriting existing files.
- If conflicting flags (`-underscore` and `-remove-underscore`) are used, the program will exit with an appropriate error message.

## Building the Tool

To build the executable, use the following command:

```
$ go build -o renamefiles main.go
```

This will create an executable named `renamefiles` that you can run in any directory.

## Running the Tool

Navigate to the directory with the files you want to rename, then run:

```
$ ./renamefiles <flags>
```

Replace `<flags>` with the desired options as explained above.

## Dependencies

This project relies on the `golang.org/x/text/unicode/norm` package for Unicode normalization.

## License

This project is licensed under the MIT License. Feel free to use and modify it as needed.
