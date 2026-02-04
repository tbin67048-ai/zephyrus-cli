package utils

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

// ReadFile reads and decrypts a file, printing its content to stdout
func ReadFile(vaultPath string, session *Session) error {
	// 1. Use FindEntry logic to navigate the nested maps
	entry, err := session.Index.FindEntry(vaultPath)
	if err != nil {
		return fmt.Errorf("could not find file in vault: %w", err)
	}

	// 2. Safety check: Ensure we aren't trying to "read" a folder
	if entry.Type == "folder" {
		return fmt.Errorf("'%s' is a directory, you can only read individual files", vaultPath)
	}

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

	// 6. Write decrypted content to stdout (no file saved)
	_, err = os.Stdout.Write(decryptedData)
	if err != nil {
		return fmt.Errorf("failed to write to stdout: %w", err)
	}

	// 7. Append newline if file doesn't end with one
	if len(decryptedData) == 0 || decryptedData[len(decryptedData)-1] != '\n' {
		_, err = os.Stdout.Write([]byte("\n"))
		if err != nil {
			return fmt.Errorf("failed to write newline: %w", err)
		}
	}

	return nil
}

// ReadSharedFile reads a shared file using a share string (username:storage_id:key)
func ReadSharedFile(shareString string) error {
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

	// 3. Fetch the encrypted file
	encryptedData, err := FetchRaw(username, storageID)
	if err != nil {
		return fmt.Errorf("failed to fetch shared file: %w", err)
	}

	// 4. Decrypt the file data
	decryptedData, err := DecryptWithKey(encryptedData, fileKey)
	if err != nil {
		return fmt.Errorf("decryption failed: %w", err)
	}

	// 5. Write decrypted content to stdout
	_, err = os.Stdout.Write(decryptedData)
	if err != nil {
		return fmt.Errorf("failed to write to stdout: %w", err)
	}

	// 6. Append newline if file doesn't end with one
	if len(decryptedData) == 0 || decryptedData[len(decryptedData)-1] != '\n' {
		_, err = os.Stdout.Write([]byte("\n"))
		if err != nil {
			return fmt.Errorf("failed to write newline: %w", err)
		}
	}

	return nil
}
