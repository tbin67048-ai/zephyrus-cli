package main

import (
	"fmt"
	"log"
	"os"

	"nexus-cli/utils"
)

func main() {
	const (
		repoURL  = "git@github.com:Auchrio/.nexus.git"
		keyPath  = ".config/key"
		fileName = "test.txt"
		fileData = "Nexus CLI Test: Pusing a file which already exists in the repo."
	)

	// 1. Read the decrypted private key from the local config
	fmt.Printf("Reading plain private key from %s...\n", keyPath)
	rawKey, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("Failed to read key file: %v. Make sure the file exists at %s", err, keyPath)
	}

	// 2. Push the file
	fmt.Printf("Pushing %s to repository...\n", fileName)

	commitMsg := "Nexus Test: pushing individual file"

	err = utils.PushFile(repoURL, rawKey, fileName, []byte(fileData), commitMsg)
	if err != nil {
		log.Fatalf("Git push failed: %v", err)
	}

	fmt.Println("--------------------------------------------------")
	fmt.Println("✔ SUCCESS: File pushed successfully!")
	fmt.Println("Check your repo: https://github.com/Auchrio/.nexus")
	fmt.Println("--------------------------------------------------")
}
