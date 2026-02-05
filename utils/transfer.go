package utils

import (
	"encoding/hex"
	"fmt"
)

// TransferVault copies all files from a source vault to a destination vault
func TransferVault(sourceUsername string, sourcePassword string, destUsername string, destPassword string) error {
	destRepoURL := fmt.Sprintf("git@github.com:%s/.zephyrus.git", destUsername)

	fmt.Printf("🔄 Starting vault transfer from %s to %s\n", sourceUsername, destUsername)

	// 1. Authenticate with source vault
	PrintProgressStep(1, 5, "Authenticating with source vault...")
	sourceSession, err := FetchSessionStateless(sourceUsername, sourcePassword)
	if err != nil {
		return fmt.Errorf("failed to authenticate with source vault: %w", err)
	}
	PrintCompletionLine("Source vault authenticated")

	// 2. Fetch destination vault SSH key (to push with)
	PrintProgressStep(2, 5, "Fetching destination vault key...")
	encryptedDestKey, err := FetchRaw(destUsername, ".config/key")
	if err != nil {
		return fmt.Errorf("destination vault not found: %w", err)
	}
	destRawKey, err := Decrypt(encryptedDestKey, destPassword)
	if err != nil {
		return fmt.Errorf("failed to decrypt destination vault key (invalid password): %w", err)
	}
	PrintCompletionLine("Destination vault authenticated")

	// 3. Prepare destination index and collect files to transfer
	PrintProgressStep(3, 5, "Scanning source vault...")
	destIndex := NewIndex()
	filesToTransfer := make(map[string][]byte)
	fileCount := 0

	// Walk through source vault and collect all files
	err = transferFilesRecursive(sourceSession, destIndex, filesToTransfer, sourceSession.Index, "", &fileCount, sourcePassword, destPassword)
	if err != nil {
		return fmt.Errorf("failed to process files: %w", err)
	}
	PrintCompletionLine(fmt.Sprintf("Found %d files to transfer", fileCount))

	if fileCount == 0 {
		return fmt.Errorf("no files found in source vault")
	}

	// 4. Encrypt destination index with destination password
	PrintProgressStep(4, 5, "Preparing transfer package...")
	destIndexBytes, err := destIndex.ToBytes(destPassword)
	if err != nil {
		return fmt.Errorf("failed to encrypt destination index: %w", err)
	}
	filesToTransfer[".config/index"] = destIndexBytes
	PrintCompletionLine("Transfer package ready")

	// 5. Push all files to destination vault
	PrintProgressStep(5, 5, "Uploading files to destination vault...")
	err = PushFilesWithAuthor(destRepoURL, destRawKey, filesToTransfer, "Zephyrus: Vault Transfer", "Zephyrus", "auchrio@proton.me")
	if err != nil {
		return fmt.Errorf("failed to upload to destination vault: %w", err)
	}
	PrintCompletionLine("Upload complete")

	fmt.Printf("✔ Successfully transferred %d files from %s to %s\n", fileCount, sourceUsername, destUsername)
	return nil
}

// transferFilesRecursive walks through the source index and transfers files
func transferFilesRecursive(
	sourceSession *Session,
	destIndex VaultIndex,
	filesToTransfer map[string][]byte,
	sourceEntries VaultIndex,
	currentPath string,
	fileCount *int,
	sourcePassword string,
	destPassword string,
) error {
	for name, entry := range sourceEntries {
		var nextPath string
		if currentPath == "" {
			nextPath = name
		} else {
			nextPath = currentPath + "/" + name
		}

		if entry.Type == "file" {
			*fileCount++
			fmt.Printf("Transferring file (%d): %s\n", *fileCount, nextPath)

			// 1. Fetch encrypted file from source
			encryptedFileData, err := FetchRaw(sourceSession.Username, entry.RealName)
			if err != nil {
				return fmt.Errorf("failed to fetch file %s: %w", nextPath, err)
			}

			// 2. Decrypt file key with source password
			encryptedKey, err := hex.DecodeString(entry.FileKey)
			if err != nil {
				return fmt.Errorf("invalid file key for %s: %w", nextPath, err)
			}
			fileKey, err := Decrypt(encryptedKey, sourcePassword)
			if err != nil {
				return fmt.Errorf("failed to decrypt file key for %s: %w", nextPath, err)
			}

			// 3. Decrypt file data with the file key
			decryptedFileData, err := DecryptWithKey(encryptedFileData, fileKey)
			if err != nil {
				return fmt.Errorf("failed to decrypt file data for %s: %w", nextPath, err)
			}

			// 4. Generate new storage name and file key for destination
			hashByteLength := sourceSession.Settings.FileHashLength / 2
			newStorageName := GenerateRandomNameWithLength(hashByteLength)
			newFileKey := GenerateFileKey()

			// 5. Encrypt new file key with destination password
			newEncryptedKey, err := Encrypt(newFileKey, destPassword)
			if err != nil {
				return fmt.Errorf("failed to encrypt new file key for %s: %w", nextPath, err)
			}
			newEncryptedKeyHex := hex.EncodeToString(newEncryptedKey)

			// 6. Encrypt file data with new file key
			newEncryptedFileData, err := EncryptWithKey(decryptedFileData, newFileKey)
			if err != nil {
				return fmt.Errorf("failed to re-encrypt file data for %s: %w", nextPath, err)
			}

			// 7. Add to destination index
			destIndex.AddFile(nextPath, newStorageName, newEncryptedKeyHex)

			// 8. Collect encrypted file for transfer
			filesToTransfer[newStorageName] = newEncryptedFileData
			fmt.Printf("  → Transferred: %s (%s)\n", nextPath, newStorageName)

		} else if entry.Type == "folder" && entry.Contents != nil {
			// Recurse into subdirectories
			err := transferFilesRecursive(sourceSession, destIndex, filesToTransfer, entry.Contents, nextPath, fileCount, sourcePassword, destPassword)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
