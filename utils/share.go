package utils

import (
	"encoding/hex"
	"fmt"
)

// ShareFile generates a share string for a file: username:storage_id:decryption_key
func ShareFile(vaultPath string, session *Session) (string, error) {
	// 1. Find the file entry in the index
	entry, err := session.Index.FindEntry(vaultPath)
	if err != nil {
		return "", fmt.Errorf("could not find file in vault: %w", err)
	}

	// 2. Ensure it's a file, not a folder
	if entry.Type == "folder" {
		return "", fmt.Errorf("'%s' is a directory, you can only share individual files", vaultPath)
	}

	// 3. Decrypt the file key from the index
	encryptedKey, err := hex.DecodeString(entry.FileKey)
	if err != nil {
		return "", fmt.Errorf("invalid file key in index: %w", err)
	}
	fileKey, err := Decrypt(encryptedKey, session.Password)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt file key: check your password")
	}

	// 4. Generate the share string: username:storage_id:key_hex
	shareString := fmt.Sprintf("%s:%s:%s", session.Username, entry.RealName, EncodeKey(fileKey))

	return shareString, nil
}
