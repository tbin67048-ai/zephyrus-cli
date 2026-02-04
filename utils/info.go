package utils

import (
	"fmt"
	"strings"
)

// VaultStats holds statistics about the vault
type VaultStats struct {
	TotalFiles   int
	TotalFolders int
	TotalSize    int64
}

// GetVaultStats calculates statistics about the vault
func GetVaultStats(session *Session) VaultStats {
	stats := VaultStats{}

	// Recursively count files and folders
	var countEntries func(Entry)
	countEntries = func(e Entry) {
		if e.Type == "file" {
			stats.TotalFiles++
		} else {
			stats.TotalFolders++
			for _, subEntry := range e.Contents {
				countEntries(subEntry)
			}
		}
	}

	// Count all entries in the index
	for _, entry := range session.Index {
		countEntries(entry)
	}

	return stats
}

// GetFileInfo retrieves detailed information about a specific file
func GetFileInfo(vaultPath string, session *Session) (map[string]interface{}, error) {
	entry, err := session.Index.FindEntry(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("could not find file in vault: %w", err)
	}

	if entry.Type == "folder" {
		return nil, fmt.Errorf("'%s' is a directory, not a file", vaultPath)
	}

	// Fetch file from remote to get size
	PrintProgressStep(1, 2, "Fetching file metadata...")
	encryptedData, err := FetchRaw(session.Username, entry.RealName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch file from remote: %w", err)
	}
	PrintCompletionLine("File metadata retrieved")

	info := map[string]interface{}{
		"name":          strings.Split(vaultPath, "/")[len(strings.Split(vaultPath, "/"))-1],
		"vaultPath":     vaultPath,
		"storageID":     entry.RealName,
		"encryptedSize": len(encryptedData),
		"fileKey":       entry.FileKey,
	}

	return info, nil
}

// PrintVaultInfo prints formatted vault information
func PrintVaultInfo(session *Session) {
	stats := GetVaultStats(session)

	fmt.Println("\n╔════════════════════════════════════════╗")
	fmt.Println("║        VAULT INFORMATION               ║")
	fmt.Println("╚════════════════════════════════════════╝")
	fmt.Printf("Username:              %s\n", session.Username)
	fmt.Printf("Total Files:           %d\n", stats.TotalFiles)
	fmt.Printf("Total Folders:         %d\n", stats.TotalFolders)
	fmt.Println("\n╔════════════════════════════════════════╗")
	fmt.Println("║      VAULT SETTINGS                    ║")
	fmt.Println("╚════════════════════════════════════════╝")
	fmt.Printf("Commit Author:         %s <%s>\n", session.Settings.CommitAuthorName, session.Settings.CommitAuthorEmail)
	fmt.Printf("Commit Message:        %s\n", session.Settings.CommitMessage)
	fmt.Printf("File Hash Length:      %d characters\n", session.Settings.FileHashLength)
	fmt.Printf("Share Hash Length:     %d characters\n", session.Settings.ShareHashLength)
	fmt.Println()
}

// PrintFileInfo prints formatted file information
func PrintFileInfo(fileInfo map[string]interface{}) {
	fmt.Println("\n╔════════════════════════════════════════╗")
	fmt.Println("║        FILE INFORMATION                ║")
	fmt.Println("╚════════════════════════════════════════╝")
	fmt.Printf("File Name:             %s\n", fileInfo["name"])
	fmt.Printf("Vault Path:            %s\n", fileInfo["vaultPath"])
	fmt.Printf("Storage ID (Hash):     %s\n", fileInfo["storageID"])
	fmt.Printf("Encrypted Size:        %d bytes\n", fileInfo["encryptedSize"])
	fmt.Printf("File Key (encrypted):  %s\n", fileInfo["fileKey"])
	fmt.Println()
}
