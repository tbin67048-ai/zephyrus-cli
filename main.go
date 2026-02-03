package main

import (
	"fmt"
	"log"
	"os"

	"nexus-cli/utils"
)

const (
	username = "Auchrio"
	repoURL  = "git@github.com:Auchrio/.nexus.git"
	keyPath  = ".config/key"
)

func main() {
	// 1. Credentials
	fmt.Print("Enter Vault Password: ")
	var password string
	fmt.Scanln(&password)

	// In a real CLI, you'd get these from os.Args or flags
	sourceFile := "test.txt"
	vaultPath := "test.txt"

	// 2. Load SSH Key
	rawKey, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("Failed to read SSH key: %v", err)
	}

	err = utils.UploadFile(sourceFile, vaultPath, password, rawKey, username, repoURL)
	if err != nil {
		log.Fatalf("Upload failed: %v", err)
	}

	fmt.Println("--------------------------------------------------")
	fmt.Println("✔ SUCCESS: Vault updated.")
	fmt.Println("--------------------------------------------------")
}
