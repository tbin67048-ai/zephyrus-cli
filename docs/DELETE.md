# delete.go Documentation

## Package utils

This module handles file and folder deletion operations in the Nexus vault, managing both the index and git repository.

### Imports

- `fmt`: String formatting and printing
- `strings`: String manipulation utilities
- `time`: Time package for commit timestamps
- `github.com/go-git/go-billy/v5/memfs`: In-memory filesystem
- `github.com/go-git/go-git/v5`: Git operations library
- `github.com/go-git/go-git/v5/config`: Git configuration
- `github.com/go-git/go-git/v5/plumbing`: Git plumbing operations
- `github.com/go-git/go-git/v5/plumbing/object`: Git object types
- `github.com/go-git/go-git/v5/plumbing/transport/ssh`: SSH transport for git
- `github.com/go-git/go-git/v5/storage/memory`: In-memory git storage
- `golang.org/x/crypto/ssh`: SSH cryptography utilities

### Functions

#### DeletePath

```go
func DeletePath(vaultPath string, session *Session) error
```

Handles both single file deletion and recursive folder deletion from the vault. This function performs the following operations:

**Parameters:**
- `vaultPath`: The path to the file or folder to delete (e.g., "folder/file.txt")
- `session`: The active session containing authentication credentials and vault index

**Return:**
- `error`: Returns an error if any operation fails (navigation, git operations, or index updates)

**Process:**

1. **Navigate to Target**: Parses the vault path and navigates through the index to locate the target file or folder
2. **Identify Files to Delete**: 
   - For files: Collects the single file's storage ID
   - For folders: Recursively collects all nested file storage IDs
3. **Git Setup**: 
   - Clones the vault git repository using SSH authentication
   - Creates an in-memory filesystem and storage
4. **Remove Files**: Physically removes all collected files from the git worktree
5. **Update Index**: Removes the target entry from the vault index
6. **Push Changes**:
   - Writes the updated index to git
   - Commits changes with author information
   - Pushes the commit back to the remote master branch

**Error Handling:**
- Returns an error if a path component is not found
- Returns an error if the target path doesn't exist in the vault
- Skips files that are already gone from the remote
- Returns an error if git operations fail

**Example Usage:**

```go
err := DeletePath("documents/report.pdf", session)
if err != nil {
    fmt.Println("Delete failed:", err)
}
```

**Notes:**

- The function uses SSH authentication with the user's SSH key stored in the session
- Supports recursive deletion of entire folders and all their contents
- All changes are automatically committed and pushed to the vault repository
- The vault index is updated to reflect the deletion
