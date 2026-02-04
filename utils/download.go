package utils

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

func DownloadFile(vaultPath string, outputPath string, session *Session) error {
	// 1. Use your custom FindEntry logic to navigate the nested maps
	entry, err := session.Index.FindEntry(vaultPath)
	if err != nil {
		return fmt.Errorf("could not find file in vault: %w", err)
	}

	// 2. Safety check: Ensure we aren't trying to "download" a folder
	if entry.Type == "folder" {
		return fmt.Errorf("'%s' is a directory, you can only download individual files", vaultPath)
	}

	fmt.Printf("Downloading %s (Storage ID: %s)...\n", vaultPath, entry.RealName)

	// 3. Fetch the encrypted hex-named file from GitHub
	encryptedData, err := FetchRaw(session.Username, entry.RealName)
	if err != nil {
		return fmt.Errorf("failed to fetch storage file from remote: %w", err)
	}

	// 4. Decrypt the file key from the index
	encryptedKey, err := hex.DecodeString(entry.FileKey)
	if err != nil {
		return fmt.Errorf("invalid file key in index: %w", err)
	}
	fileKey, err := Decrypt(encryptedKey, session.Password)
	if err != nil {
		return fmt.Errorf("failed to decrypt file key: check your password")
	}

	// 5. Decrypt the file data with the file key
	decryptedData, err := DecryptWithKey(encryptedData, fileKey)
	if err != nil {
		return fmt.Errorf("decryption failed: %w", err)
	}

	// 6. Save to the local output path
	return os.WriteFile(outputPath, decryptedData, 0644)
}

// DownloadSharedFile downloads a file using a share string (username:storage_id:key)
func DownloadSharedFile(shareString string, outputPath string) error {
	// 1. Parse the share string
	parts := strings.Split(shareString, ":")
	if len(parts) != 3 {
		return fmt.Errorf("invalid share string format, expected 'username:storage_id:key'")
	}

	username := parts[0]
	storageID := parts[1]
	keyHex := parts[2]

	// 2. Decode the file key
	fileKey, err := DecodeKey(keyHex)
	if err != nil {
		return fmt.Errorf("invalid file key in share string: %w", err)
	}

	fmt.Printf("Downloading shared file from %s (Storage ID: %s)...\n", username, storageID)

	// 3. Fetch the encrypted file from GitHub
	encryptedData, err := FetchRaw(username, storageID)
	if err != nil {
		return fmt.Errorf("failed to fetch shared file from remote: %w", err)
	}

	// 4. Decrypt with the provided file key
	decryptedData, err := DecryptWithKey(encryptedData, fileKey)
	if err != nil {
		return fmt.Errorf("decryption failed: invalid share key")
	}

	// 5. Save to the local output path
	return os.WriteFile(outputPath, decryptedData, 0644)
}
