# upload.go Documentation

## Package utils

This module handles uploading files to the vault, including encryption, index updates, and remote repository synchronization.

### Imports

- `fmt`: String formatting and printing
- `os`: Operating system file operations

### Functions

#### UploadFile

```go
func UploadFile(sourcePath string, vaultPath string, session *Session) error
```

Uploads a file from the local filesystem to the vault, encrypting it and updating the vault index.

**Parameters:**
- `sourcePath`: The path to the file on the local filesystem (e.g., "./document.pdf")
- `vaultPath`: The destination path in the vault (e.g., "documents/report.pdf")
- `session`: The active session containing authentication credentials, encryption password, and vault index

**Return:**
- `error`: Returns error if file reading, encryption, index update, or push operations fail

**Upload Process:**

1. **Construct Repository URL**: Builds the git repository URL from session username

2. **Read Source File**: Reads the entire file from the local filesystem
   - Returns error if file doesn't exist or cannot be read

3. **Determine Storage Name**:
   - **Existing Files**: Checks if file already exists in vault index
     - Uses existing storage ID (RealName)
     - Displays "Updating existing file" message
   - **New Files**: Generates new random hex-based storage ID
     - Adds entry to vault index with `AddFile()`
     - Displays "Uploading new file" message

4. **Encrypt File Data**:
   - Encrypts file using the session password
   - Produces encrypted bytes ready for storage

5. **Encrypt Updated Index**:
   - Serializes and encrypts the updated vault index
   - Creates `.config/index` entry

6. **Prepare Push Package**:
   - Creates map containing:
     - The encrypted file (storage ID as key)
     - The encrypted index (`.config/index` as key)

7. **Push to Repository**:
   - Calls `PushFiles()` to push both files to GitHub
   - Uses SSH authentication from session
   - Commit message: "Nexus: Updated Vault"

**File Handling:**

**New Files:**
```
Local: ./document.pdf
Vault: documents/report.pdf
Storage: a3f2e1c9d4b6f8e2 (hex name, content = encrypted)
```

**Updating Existing Files:**
```
Vault: documents/report.pdf
Storage: (existing hex name - reused)
Action: File content replaced, index unchanged
```

**Error Handling:**
- Returns error if source file cannot be read
- Returns error if encryption fails
- Returns error if git push fails
- Returns error if index serialization fails

**Example Usage:**

```go
// Upload a new file
err := UploadFile("./quarterly_report.pdf", "documents/Q1.pdf", session)
if err != nil {
    fmt.Println("Upload failed:", err)
}

// Update an existing file
err = UploadFile("./updated_report.pdf", "documents/Q1.pdf", session)
if err != nil {
    fmt.Println("Update failed:", err)
}

// Upload to nested folders (creates folders if needed)
err = UploadFile("./archive.zip", "backups/2024/archive.zip", session)
if err != nil {
    fmt.Println("Upload failed:", err)
}
```

**Output Examples:**

Uploading new file:
```
Uploading new file: documents/report.pdf as a3f2e1c9d4b6f8e2
```

Updating existing file:
```
Updating existing file: documents/report.pdf (a3f2e1c9d4b6f8e2)
```

### Implementation Details

The upload process:
- Checks local index cache (not remote) for existing file
- Creates intermediate folders automatically
- Uses stateless push operation (doesn't redownload entire repository)
- Updates both file content and vault index in single push

### Security Considerations

- Files are encrypted with AES-256-GCM before upload
- Storage IDs are randomly generated (original filenames not exposed)
- Index is encrypted and stored in `.config/index`
- Encryption key is derived from session password

### File Size Considerations

- No explicit file size limits in the code
- Practical limits depend on GitHub's API and SSH transfer limits
- Large files should be zipped before upload

### Notes

- If an upload fails mid-process, the local session index may be updated but the push might fail
- The vault index is updated before pushing (optimistic approach)
- Consider saving the session after successful upload
- Intermediate folders are created automatically (no need to pre-create structure)
- Files with the same vault path but different source paths will overwrite

### Typical Upload Workflow

```go
// After connecting to vault
err := UploadFile("./my_document.pdf", "documents/my_document.pdf", session)
if err != nil {
    fmt.Println("Upload failed:", err)
    return
}

// File is now encrypted and stored in vault
// Index is updated and pushed to GitHub
// Both files and index are synced with remote
```
