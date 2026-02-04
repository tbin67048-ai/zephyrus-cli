# index.go Documentation

## Package utils

This module manages the vault index structure, which is a hierarchical mapping of file and folder names to their encrypted storage IDs.

### Types

#### Entry

```go
type Entry struct {
    Type     string           `json:"type"`
    RealName string           `json:"realName,omitempty"`
    Contents map[string]Entry `json:"contents,omitempty"`
}
```

Represents either a file or folder in the Nexus vault.

**Fields:**
- **Type**: Either "file" or "folder"
- **RealName**: The encrypted storage ID (hex name) of the file; omitted for folders
- **Contents**: A map of child entries; used only for folders

#### VaultIndex

```go
type VaultIndex map[string]Entry
```

The top-level structure representing the entire vault. Maps entry names to their `Entry` objects at the root level.

### Functions

#### NewIndex

```go
func NewIndex() VaultIndex
```

Creates an empty vault index.

**Return:**
- `VaultIndex`: A new, empty index ready for use

**Example:**
```go
index := NewIndex()
```

#### FromBytes

```go
func FromBytes(data []byte, password string) (VaultIndex, error)
```

Decrypts and parses an encrypted index from bytes.

**Parameters:**
- `data`: The encrypted index data
- `password`: The password used to decrypt the index

**Return:**
- `VaultIndex`: The decrypted and parsed index
- `error`: Returns error if decryption or JSON parsing fails

**Process:**
1. Decrypts the data using the password
2. Unmarshals the JSON into a VaultIndex

#### ToBytes

```go
func (vi VaultIndex) ToBytes(password string) ([]byte, error)
```

Serializes and encrypts the index to bytes for storage.

**Parameters:**
- `password`: The password used to encrypt the index

**Return:**
- `[]byte`: The encrypted index data
- `error`: Returns error if JSON marshaling or encryption fails

**Process:**
1. Marshals the index to indented JSON
2. Encrypts the JSON using the password

#### FindEntry

```go
func (vi VaultIndex) FindEntry(path string) (*Entry, error)
```

Navigates the hierarchical index to find an entry at a specified path.

**Parameters:**
- `path`: The path to the entry (e.g., "documents/reports/quarterly.pdf")

**Return:**
- `*Entry`: Pointer to the found entry
- `error`: Returns error if path component is not found or path is invalid

**Error Cases:**
- Path component not found
- Intermediate path component is a file (not a folder)
- Invalid path format

**Example:**
```go
entry, err := index.FindEntry("documents/report.pdf")
if err != nil {
    log.Fatal(err)
}
fmt.Println(entry.RealName) // "a3f2e1c9d4b6f8e2"
```

#### AddFile

```go
func (vi VaultIndex) AddFile(path string, realName string)
```

Inserts a new file entry into the index, creating intermediate folders if necessary.

**Parameters:**
- `path`: The vault path for the file (e.g., "documents/reports/Q1.pdf")
- `realName`: The encrypted storage ID to associate with the file

**Behavior:**
- Automatically creates folder entries as needed
- If intermediate paths don't exist, they are created with type "folder"
- The final entry is created with type "file"

**Example:**
```go
index.AddFile("documents/reports/Q1.pdf", "a3f2e1c9d4b6f8e2")
// Creates "documents" folder and "reports" subfolder if they don't exist
// Then creates "Q1.pdf" file entry with the given RealName
```

#### PrintDebug

```go
func (vi VaultIndex) PrintDebug()
```

Prints the entire index structure to the console in formatted JSON.

**Usage:**
- Debug tool to view the current vault structure
- Prints "Index is currently empty." if index is empty

**Example:**
```go
index.PrintDebug()
// Output:
// {
//   "documents": {
//     "type": "folder",
//     "contents": {
//       "report.pdf": {
//         "type": "file",
//         "realName": "a3f2e1c9d4b6f8e2"
//       }
//     }
//   }
// }
```

### Example Index Structure (JSON)

```json
{
  "documents": {
    "type": "folder",
    "contents": {
      "report.pdf": {
        "type": "file",
        "realName": "a3f2e1c9d4b6f8e2"
      },
      "archive": {
        "type": "folder",
        "contents": {
          "old_report.pdf": {
            "type": "file",
            "realName": "f1e2d3c4b5a6f7e8"
          }
        }
      }
    }
  },
  "images": {
    "type": "folder",
    "contents": {}
  }
}
```

### Notes

- The index is encrypted and stored in `.config/index` in the vault repository
- File names are obfuscated; the original names are stored in the index, while files are stored with hex storage IDs
- Folders are represented with `type: "folder"` and empty `Contents` initially
- The index is human-readable (before encryption) for debugging purposes
