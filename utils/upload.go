package utils

import (
	"fmt"
	"os"
)

func UploadFile(sourcePath string, vaultPath string, session *Session) error {
	repoURL := fmt.Sprintf("git@github.com:%s/.nexus.git", session.Username)

	// 1. Read source
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	// 2. Determine Storage Name (Check local index cache)
	var realName string
	if entry, exists := session.Index[vaultPath]; exists {
		realName = entry.RealName
		fmt.Printf("Updating existing file: %s (%s)\n", vaultPath, realName)
	} else {
		realName = GenerateRandomName()
		session.Index.AddFile(vaultPath, realName)
		fmt.Printf("Uploading new file: %s as %s\n", vaultPath, realName)
	}

	// 3. Encrypt file data
	encryptedData, err := Encrypt(data, session.Password)
	if err != nil {
		return err
	}

	// 4. Encrypt updated index
	indexBytes, err := session.Index.ToBytes(session.Password)
	if err != nil {
		return err
	}

	// 5. Push to Git
	filesToPush := map[string][]byte{
		realName:        encryptedData,
		".config/index": indexBytes,
	}

	err = PushFiles(repoURL, session.RawKey, filesToPush, "Nexus: Updated Vault")
	if err != nil {
		return err
	}

	// 6. Save updated index to local session to bypass cache
	return nil
}
