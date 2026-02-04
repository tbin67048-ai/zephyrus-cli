# download.go Documentation

## Package utils

This module handles downloading files from the Nexus vault, including retrieval, decryption, and local storage operations.

### Imports

- `fmt`: String formatting and printing
- `os`: Operating system file operations

### Functions

#### DownloadFile

```go
func DownloadFile(vaultPath string, outputPath string, session *Session) error
```

Downloads a file from the vault, decrypts it, and saves it to the local filesystem.

**Parameters:**
- `vaultPath`: The path to the file in the vault (e.g., "documents/report.pdf")
- `outputPath`: The local filesystem path where the file will be saved
- `session`: The active session containing authentication credentials and vault index

**Return:**
- `error`: Returns an error if any operation fails (file lookup, fetching, decryption, or writing)

**Process:**

1. **Find Entry**: Uses the vault index to locate the file entry at the specified vault path
2. **Validate Type**: Ensures the target is a file, not a folder (folders cannot be downloaded directly)
3. **Fetch Encrypted File**: Retrieves the encrypted file from the remote repository using its storage ID
4. **Decrypt**: Decrypts the file using the session password
5. **Save Locally**: Writes the decrypted file to the specified output path on the local filesystem

**Error Handling:**
- Returns an error if the file is not found in the vault
- Returns an error if attempting to download a folder
- Returns an error if the remote file fetch fails
- Returns an error if decryption fails (usually indicates invalid password)
- Returns an error if local file write fails

**Example Usage:**

```go
err := DownloadFile("documents/report.pdf", "./report.pdf", session)
if err != nil {
    fmt.Println("Download failed:", err)
}
```

**Notes:**

- Files are stored in the vault with hex-based storage IDs (RealName) for privacy
- The actual filename stored remotely does not reveal the original file name
- Downloaded files are created with permission mode 0644 (read/write for owner, read-only for others)
- Decryption failures typically indicate an incorrect password or corrupted file data
