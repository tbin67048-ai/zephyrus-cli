# search.go Documentation

## Package utils

This module provides search functionality to find files and folders in the vault by name or path.

### Imports

- `fmt`: String formatting and printing
- `os`: Operating system interfaces
- `strings`: String manipulation utilities
- `text/tabwriter`: Formatted tabular output

### Functions

#### SearchFiles

```go
func SearchFiles(session *Session, query string) error
```

Searches the vault for files and folders matching a query string and displays results in formatted tabular output.

**Parameters:**
- `session`: The active session containing the vault index
- `query`: The search query string (case-insensitive partial match)

**Return:**
- `error`: Always returns nil; errors from search don't prevent output

**Behavior:**

1. **Parse Query**: Converts query to lowercase for case-insensitive search
2. **Recursive Search**: Walks the entire vault index tree
3. **Match Criteria**: Matches if query string is contained in:
   - The entry name
   - The full vault path
4. **Display Results**: Shows matching entries in formatted table with three columns:
   - **VAULT PATH**: The full path from root (e.g., "documents/reports/Q1.pdf")
   - **TYPE**: "[FILE]" or "[DIR]"
   - **STORAGE ID**: The hex-encoded storage ID for files, "-" for folders

**Search Example:**

Query: "report"

Results:
```
VAULT PATH              TYPE    STORAGE ID
----------              ----    ----------
documents/report.pdf    [FILE]  a3f2e1c9d4b6f8e2
archive/2024_report     [DIR]   -
reports/Q1_report.pdf   [FILE]  b4c6f7e2d3a1f8e9
```

**Search Behavior:**

- **Case-Insensitive**: "REPORT", "Report", "report" all match the same files
- **Partial Match**: "rep" matches "report", "2024_report", etc.
- **Full Path Matching**: Searching for "documents/report" matches files in that path
- **Folder Matching**: Can find both files and folders

**Example Usage:**

```go
// Search for all PDFs
err := SearchFiles(session, "pdf")
if err != nil {
    log.Fatal(err)
}

// Search for reports
err = SearchFiles(session, "report")

// Search by partial path
err = SearchFiles(session, "documents/2024")

// Search with specific extension
err = SearchFiles(session, ".xlsx")
```

**Output When No Matches Found:**

```
VAULT PATH    TYPE    STORAGE ID
----------    ----    ----------
No matches found.
```

### Implementation Details

The search uses a recursive helper function `walk()` that:
- Iterates through entries at current level
- Constructs full paths as it traverses
- Checks if full path contains query (case-insensitive)
- Recursively processes folder contents
- Tracks whether any matches were found

### Notes

- Search is performed in-memory on the vault index (fast)
- Search results show full vault paths, not just filenames
- Folders are displayed with "/" appended to their names in results
- Storage IDs are the actual hex file names used in the repository
- Query matching is substring-based (not regex)
- Case-insensitive matching helps find files regardless of naming conventions
- The order of results depends on Go's map iteration (not guaranteed)

### Common Search Patterns

```go
// Find all files of a type
SearchFiles(session, ".pdf")    // All PDFs
SearchFiles(session, ".xlsx")   // All Excel files

// Find files in a specific folder
SearchFiles(session, "documents/invoices")

// Find files by partial name
SearchFiles(session, "2024")    // All 2024 files
SearchFiles(session, "backup")  // All backups
```
