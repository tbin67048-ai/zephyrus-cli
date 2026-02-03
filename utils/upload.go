package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// uploadFile now takes rawKey, username, and repoURL as arguments.
func UploadFile(sourcePath string, vaultPath string, password string, rawKey []byte, username string, repoURL string) error {

	// 1. Fetch and Check Index
	fmt.Println("Fetching current index...")
	var index VaultIndex
	var realName string

	// FetchRaw uses the username provided
	rawIndex, err := FetchRaw(username, ".config/index")
	if err != nil {
		fmt.Println("No existing index found, creating a new one.")
		index = NewIndex()
	} else {
		index, err = FromBytes(rawIndex, password)
		if err != nil {
			return fmt.Errorf("failed to decrypt index: %v", err)
		}

		// Check if file exists to reuse the 16-char filename (RealName)
		if entry, exists := index[vaultPath]; exists {
			fmt.Printf("⚠ CONFLICT: '%s' already exists (ID: %s).\n", vaultPath, entry.RealName)
			fmt.Print("Do you want to overwrite it? (y/N): ")

			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.ToLower(strings.TrimSpace(response))

			if response != "y" && response != "yes" {
				return fmt.Errorf("upload cancelled by user")
			}

			// Reuse the existing ID to perform a true overwrite in Git
			realName = entry.RealName
			fmt.Printf("Proceeding with overwrite using existing ID: %s\n", realName)
		}
	}

	// 2. Prepare File Data
	content, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read source file: %v", err)
	}

	encryptedFile, err := Encrypt(content, password)
	if err != nil {
		return fmt.Errorf("encryption failed: %v", err)
	}

	// If it's a new file (not an overwrite), generate a new ID
	if realName == "" {
		realName = GenerateRandomName()
		index.AddFile(vaultPath, realName)
	}

	encryptedIndex, err := index.ToBytes(password)
	if err != nil {
		return fmt.Errorf("failed to encrypt index: %v", err)
	}

	// 3. Atomic Push
	filesToPush := map[string][]byte{
		realName:        encryptedFile,
		".config/index": encryptedIndex,
	}

	fmt.Printf("Pushing %s (as %s)...\n", vaultPath, realName)

	// PushFiles uses the rawKey passed into the function
	return PushFiles(repoURL, rawKey, filesToPush, "Nexus: Updated "+vaultPath)
}
