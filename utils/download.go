package utils

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func DownloadFile(vaultPath string, outputPath string, session *Session) error {
	// 1. Use your custom FindEntry logic to navigate the nested maps
	PrintProgressStep(1, 5, "Locating file in vault...")
	entry, err := session.Index.FindEntry(vaultPath)
	if err != nil {
		return fmt.Errorf("could not find file in vault: %w", err)
	}

	// 2. Safety check: Ensure we aren't trying to "download" a folder
	if entry.Type == "folder" {
		return fmt.Errorf("'%s' is a directory, you can only download individual files", vaultPath)
	}
	PrintCompletionLine("File located: " + entry.RealName)

	fmt.Printf("Downloading %s (Storage ID: %s)...\n", vaultPath, entry.RealName)

	// 3. Fetch the encrypted hex-named file from GitHub
	PrintProgressStep(2, 5, "Fetching encrypted file from GitHub...")
	encryptedData, err := FetchRaw(session.Username, entry.RealName)
	if err != nil {
		return fmt.Errorf("failed to fetch storage file from remote: %w", err)
	}
	PrintCompletionLine("File fetched from GitHub")

	// 4. Decrypt the file key from the index
	PrintProgressStep(3, 5, "Decrypting file key...")
	encryptedKey, err := hex.DecodeString(entry.FileKey)
	if err != nil {
		return fmt.Errorf("invalid file key in index: %w", err)
	}
	fileKey, err := Decrypt(encryptedKey, session.Password)
	if err != nil {
		return fmt.Errorf("failed to decrypt file key: check your password")
	}
	PrintCompletionLine("File key decrypted")

	// 5. Decrypt the file data with the file key
	PrintProgressStep(4, 5, "Decrypting file contents...")
	time.Sleep(time.Millisecond * 100) // Simulate work for visibility
	decryptedData, err := DecryptWithKey(encryptedData, fileKey)
	if err != nil {
		return fmt.Errorf("decryption failed: %w", err)
	}
	PrintCompletionLine("File contents decrypted")

	// 6. Save to the local output path
	PrintProgressStep(5, 5, "Saving file to "+outputPath+"...")
	err = os.WriteFile(outputPath, decryptedData, 0644)
	if err != nil {
		return err
	}
	PrintCompletionLine("File saved successfully")
	return nil
}

// DownloadDirectory downloads an entire directory recursively from the vault
func DownloadDirectory(vaultPath string, outputPath string, session *Session) error {
	// 1. Verify the path is a directory
	PrintProgressStep(1, 3, "Locating directory in vault...")
	entry, err := session.Index.FindEntry(vaultPath)
	if err != nil {
		return fmt.Errorf("could not find directory in vault: %w", err)
	}

	// 2. Safety check: Ensure we're downloading a folder
	if entry.Type != "folder" {
		return fmt.Errorf("'%s' is a file, not a directory. Use download command for files", vaultPath)
	}
	PrintCompletionLine("Directory located")

	fmt.Printf("Downloading directory from vault: %s\n", vaultPath)

	// 3. Create output directory if it doesn't exist
	err = os.MkdirAll(outputPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	fileCount := 0

	// 4. Recursively download all files in the directory
	var downloadFiles func(currentEntry Entry, currentVaultPath string, currentLocalPath string) error
	downloadFiles = func(currentEntry Entry, currentVaultPath string, currentLocalPath string) error {
		// Process all entries in the current folder
		for name, subEntry := range currentEntry.Contents {
			var nextVaultPath string
			if currentVaultPath == "" {
				nextVaultPath = name
			} else {
				nextVaultPath = currentVaultPath + "/" + name
			}

			nextLocalPath := filepath.Join(currentLocalPath, name)

			if subEntry.Type == "file" {
				fileCount++
				fmt.Printf("Downloading file (%d): %s\n", fileCount, name)

				// 5. Fetch the encrypted file from GitHub
				encryptedData, err := FetchRaw(session.Username, subEntry.RealName)
				if err != nil {
					return fmt.Errorf("failed to fetch file %s: %w", nextVaultPath, err)
				}

				// 6. Decrypt the file key from the index
				encryptedKey, err := hex.DecodeString(subEntry.FileKey)
				if err != nil {
					return fmt.Errorf("invalid file key in index for %s: %w", nextVaultPath, err)
				}
				fileKey, err := Decrypt(encryptedKey, session.Password)
				if err != nil {
					return fmt.Errorf("failed to decrypt file key for %s: %w", nextVaultPath, err)
				}

				// 7. Decrypt the file data with the file key
				decryptedData, err := DecryptWithKey(encryptedData, fileKey)
				if err != nil {
					return fmt.Errorf("decryption failed for %s: %w", nextVaultPath, err)
				}

				// 8. Save to local path
				err = os.WriteFile(nextLocalPath, decryptedData, 0644)
				if err != nil {
					return fmt.Errorf("failed to save file %s: %w", nextLocalPath, err)
				}

				fmt.Printf("  → Saved: %s\n", nextLocalPath)

			} else if subEntry.Type == "folder" {
				// Create subdirectory
				err := os.MkdirAll(nextLocalPath, 0755)
				if err != nil {
					return fmt.Errorf("failed to create directory %s: %w", nextLocalPath, err)
				}

				// Recursively download contents
				err = downloadFiles(subEntry, nextVaultPath, nextLocalPath)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}

	// Start recursive download from the directory entry
	err = downloadFiles(*entry, vaultPath, outputPath)
	if err != nil {
		return err
	}

	if fileCount == 0 {
		return fmt.Errorf("no files found in directory: %s", vaultPath)
	}

	PrintProgressStep(2, 3, "Finalizing download...")
	PrintCompletionLine("Decryption complete")

	PrintProgressStep(3, 3, "Writing files to disk...")
	PrintCompletionLine("Files written successfully")

	fmt.Printf("✔ Successfully downloaded %d files from directory\n", fileCount)
	return nil
}

// DownloadSharedFile downloads a file using a share string (username:reference:sharepassword:base64filename)
func DownloadSharedFile(shareString string, outputPath string) error {
	// 1. Parse the share string (supports both old 3-part and new 4-part formats)
	parts := strings.Split(shareString, ":")
	if len(parts) < 3 || len(parts) > 4 {
		return fmt.Errorf("invalid share string format, expected 'username:reference:sharepassword' or 'username:reference:sharepassword:base64filename'")
	}

	username := parts[0]
	reference := parts[1]
	sharePassword := parts[2]
	var filename string

	// Decode filename if provided in share string
	if len(parts) == 4 {
		decoded, err := base64.StdEncoding.DecodeString(parts[3])
		if err != nil {
			return fmt.Errorf("invalid filename encoding in share string: %w", err)
		}
		filename = string(decoded)
	}

	// Use provided outputPath, or construct from filename if available
	finalOutputPath := outputPath
	if filename != "" && outputPath == "" {
		finalOutputPath = filename
	}

	if filename != "" {
		fmt.Printf("Downloading '%s' from %s (Reference: %s)...\n", filename, username, reference)
	} else {
		fmt.Printf("Downloading shared file from %s (Reference: %s)...\n", username, reference)
	}

	// 2. Fetch the share pointer from the /shared/ folder
	sharedPath := fmt.Sprintf("shared/%s", reference)
	pointerData, err := FetchRaw(username, sharedPath)
	if err != nil {
		return fmt.Errorf("failed to fetch share pointer from remote: %w", err)
	}

	// 3. Decrypt the pointer with the share password to get storage ID and file key
	decryptedPointer, err := Decrypt(pointerData, sharePassword)
	if err != nil {
		return fmt.Errorf("decryption failed: invalid share password")
	}

	// 4. Parse the pointer JSON to get storage ID and encrypted file key
	var pointerMap map[string]string
	err = json.Unmarshal(decryptedPointer, &pointerMap)
	if err != nil {
		return fmt.Errorf("invalid share pointer format: %w", err)
	}

	storageID, ok := pointerMap["storageID"]
	if !ok {
		return fmt.Errorf("share pointer missing storageID")
	}
	fileKeyHex, ok := pointerMap["fileKey"]
	if !ok {
		return fmt.Errorf("share pointer missing fileKey")
	}

	// 5. Fetch the actual encrypted file from main storage
	encryptedFileData, err := FetchRaw(username, storageID)
	if err != nil {
		return fmt.Errorf("failed to fetch file from remote: %w", err)
	}

	// 6. Decrypt the file key (it's encrypted with the share password for this transfer)
	encryptedKeyBytes, err := hex.DecodeString(fileKeyHex)
	if err != nil {
		return fmt.Errorf("invalid file key encoding: %w", err)
	}
	fileKey, err := Decrypt(encryptedKeyBytes, sharePassword)
	if err != nil {
		return fmt.Errorf("failed to decrypt file key: %w", err)
	}

	// 7. Decrypt the file content with the file key
	decryptedData, err := DecryptWithKey(encryptedFileData, fileKey)
	if err != nil {
		return fmt.Errorf("file decryption failed: %w", err)
	}

	// 8. Save to the local output path
	return os.WriteFile(finalOutputPath, decryptedData, 0644)
}
