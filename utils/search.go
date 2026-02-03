package utils

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func SearchFiles(session *Session, query string) error {
	fmt.Printf("Searching vault for: \"%s\"\n", query)

	if len(session.Index) == 0 {
		fmt.Println("Vault is empty.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "VAULT PATH\tSTORAGE ID (HEX)")
	fmt.Fprintln(w, "----------\t----------------")

	foundCount := 0
	lowerQuery := strings.ToLower(query)

	for vPath, entry := range session.Index {
		if strings.Contains(strings.ToLower(vPath), lowerQuery) {
			fmt.Fprintf(w, "%s\t%s\n", vPath, entry.RealName)
			foundCount++
		}
	}

	w.Flush()
	if foundCount == 0 {
		fmt.Printf("\nNo files found matching \"%s\".\n", query)
	}
	return nil
}
