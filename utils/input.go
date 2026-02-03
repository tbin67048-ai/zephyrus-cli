package utils

import (
	"fmt"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// GetPassword prompts the user for a password without echoing input to the terminal
func GetPassword(prompt string) (string, error) {
	fmt.Print(prompt)

	// syscall.Stdin is the file descriptor for standard input
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	// ReadPassword doesn't capture the newline character, so we print one manually
	fmt.Println()

	return strings.TrimSpace(string(bytePassword)), nil
}
