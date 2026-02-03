package utils

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

// SearchFiles filters the vault index based on a query string
func SearchFiles(username string, password string, query string) error {
	fmt.Printf("Searching vault '%s' for: \"%s\"\n", username, query)

	// 1. Fetch the encrypted index
	rawIndex, err := FetchRaw(username, ".config/index")
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			fmt.Println("Vault is empty. Nothing to search.")
			return nil
		}
		return fmt.Errorf("failed to fetch index: %w", err)
	}

	// 2. Decrypt the index
	index, err := FromBytes(rawIndex, password)
	if err != nil {
		return fmt.Errorf("failed to decrypt index: %w", err)
	}

	// 3. Filter and Print
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "VAULT PATH\tSTORAGE ID (HEX)")
	fmt.Fprintln(w, "----------\t----------------")

	foundCount := 0
	lowerQuery := strings.ToLower(query)

	for vPath, entry := range index {
		// Case-insensitive match
		if strings.Contains(strings.ToLower(vPath), lowerQuery) {
			fmt.Fprintf(w, "%s\t%s\n", vPath, entry.RealName)
			foundCount++
		}
	}

	w.Flush()

	if foundCount == 0 {
		fmt.Printf("\nNo files found matching \"%s\".\n", query)
	} else {
		fmt.Printf("\nFound %d match(es).\n", foundCount)
	}

	return nil
}
