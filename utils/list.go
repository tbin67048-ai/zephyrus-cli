package utils

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func ListFiles(session *Session) error {
	if len(session.Index) == 0 {
		fmt.Println("--------------------------------------------------")
		fmt.Println("ℹ Your vault is currently empty.")
		fmt.Println("--------------------------------------------------")
		return nil
	}

	fmt.Printf("Vault Index for %s (Local Cache):\n\n", session.Username)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "VAULT PATH\tSTORAGE ID (HEX)")
	fmt.Fprintln(w, "----------\t----------------")

	for vPath, entry := range session.Index {
		fmt.Fprintf(w, "%s\t%s\n", vPath, entry.RealName)
	}

	w.Flush()
	fmt.Printf("\nTotal files: %d\n", len(session.Index))
	return nil
}
