# purge.go Documentation

## Package utils

This module provides functionality to completely wipe the vault by force-pushing an empty repository state to GitHub.

### Imports

- `fmt`: String formatting and printing
- `time`: Time package for commit timestamps
- `github.com/go-git/go-billy/v5/memfs`: In-memory filesystem
- `github.com/go-git/go-git/v5`: Git operations library
- `github.com/go-git/go-git/v5/config`: Git configuration
- `github.com/go-git/go-git/v5/plumbing/object`: Git object types
- `github.com/go-git/go-git/v5/plumbing/transport/ssh`: SSH transport for git
- `github.com/go-git/go-git/v5/storage/memory`: In-memory git storage
- `golang.org/x/crypto/ssh`: SSH cryptography utilities

### Functions

#### PurgeVault

```go
func PurgeVault(session *Session) error
```

Completely wipes the remote vault by force-pushing an empty commit history to GitHub. This is a destructive operation that removes all files and history from the vault.

**Parameters:**
- `session`: The active session containing user credentials and vault information

**Return:**
- `error`: Returns error if any operation fails (private key loading, commit, or push)

**Process:**

1. **Construct Repo URL**: Builds the git repository URL from session username
2. **Create Fresh Git Environment**: Initializes a new, empty git repository in memory
3. **Setup SSH Authentication**: Loads the private key from the session
4. **Create Wipe Commit**: Creates an empty commit with message "Nexus: Updated Vault"
5. **Connect Remote**: Adds the GitHub repository as the remote origin
6. **Force Push**: Force-pushes the empty state to GitHub, overwriting all history
7. **Clear Memory Cache**: Resets the session index to an empty vault

**Error Handling:**
- Returns error if private key loading fails
- Returns error if commit creation fails
- Returns error if push fails

**Important Notes:**

⚠️ **DESTRUCTIVE OPERATION**: This function:
- Permanently deletes all files in the vault
- Erases all commit history
- Cannot be undone

Use with extreme caution!

**Security Considerations:**
- Requires SSH authentication with the vault's private key
- Force push requires write access to the repository
- Completely clears the in-memory index (`session.Index`)

**Example Usage:**

```go
// WARNING: This is a destructive operation!
err := PurgeVault(session)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Vault purged successfully")
```

### Typical Workflow

```go
// After confirming user intent
fmt.Println("WARNING: This will permanently delete all vault data!")
fmt.Print("Type 'DELETE' to confirm: ")
var confirm string
fmt.Scanln(&confirm)

if confirm == "DELETE" {
    err := PurgeVault(session)
    if err != nil {
        fmt.Println("Purge failed:", err)
    } else {
        fmt.Println("Vault has been purged")
    }
}
```

### Related Operations

- **DisconnectAndDelete**: Removes local config file
- **DeletePath**: Removes specific files without purging entire vault
- **Disconnect**: Clears memory cache without affecting remote

### Notes

- Uses shallow clone for efficiency when checking repo state
- Ignores SSH host key verification (for automation)
- Empty commit is allowed via `AllowEmptyCommits: true`
- Force push (`Force: true`) overwrites remote history
- After purge, the session index is reset to `NewIndex()` (empty)
