# list.go Documentation

## Package utils

This module provides functionality to list files and folders in the vault with formatted output.

### Imports

- `fmt`: String formatting and printing
- `os`: Operating system interfaces
- `text/tabwriter`: Formatted tabular output

### Functions

#### ListFiles

```go
func ListFiles(session *Session, folderPath string) error
```

Lists all files and folders in a specified vault directory with formatted tabular output.

**Parameters:**
- `session`: The active session containing the vault index
- `folderPath`: The path to the folder to list (e.g., "documents" or "documents/archive")
  - Empty string or "/" lists the root of the vault
  - Must point to a folder, not a file

**Return:**
- `error`: Returns error if the path is invalid or doesn't point to a folder

**Behavior:**

1. **Navigate to Target**: If a path is provided, finds the corresponding folder entry in the index
2. **Validate**: Ensures the target is a folder, not a file
3. **Display**: Shows a formatted table with three columns:
   - **NAME**: File or folder name (folders have "/" appended)
   - **TYPE**: "[FILE]" or "[DIR]"
   - **STORAGE ID**: The hex-encoded storage ID (RealName) for files, "-" for folders

**Table Format Example:**

```
NAME                TYPE    STORAGE ID
----                ----    ----------
document.pdf        [FILE]  a3f2e1c9d4b6f8e2
archive/            [DIR]   -
image.png           [FILE]  f1e2d3c4b5a6f7e8
```

**Error Handling:**
- Returns error if the path is not found in the vault
- Returns error if the specified path is a file instead of a folder
- Displays "Directory is empty." if the folder has no entries

**Example Usage:**

```go
// List root of vault
err := ListFiles(session, "")
if err != nil {
    log.Fatal(err)
}

// List a subfolder
err = ListFiles(session, "documents/archive")
if err != nil {
    log.Fatal(err)
}
```

**Output Examples:**

List root directory:
```
NAME            TYPE    STORAGE ID
----            ----    ----------
documents/      [DIR]   -
images/         [DIR]   -
config.json     [FILE]  a3f2e1c9d4b6f8e2
```

List subfolder:
```
NAME            TYPE    STORAGE ID
----            ----    ----------
Q1_report.pdf   [FILE]  b4c6f7e2d3a1f8e9
Q2_report.pdf   [FILE]  c5d7a1f3e4b2f9e0
archive/        [DIR]   -
```

### Notes

- The table uses tab-separated columns with 3-space padding
- Folders are displayed with a trailing "/" in the NAME column
- Storage IDs (RealName) are the actual hex file names stored in the vault
- The list is displayed in no guaranteed order (depends on Go's map iteration)
- Empty directories show "Directory is empty." message
