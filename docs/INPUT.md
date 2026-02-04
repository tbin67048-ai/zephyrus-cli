# input.go Documentation

## Package utils

This module provides secure user input handling, specifically for password entry without echo.

### Imports

- `fmt`: String formatting and printing
- `strings`: String manipulation utilities
- `syscall`: System call operations
- `golang.org/x/term`: Terminal I/O utilities

### Functions

#### GetPassword

```go
func GetPassword(prompt string) (string, error)
```

Prompts the user for a password without echoing input to the terminal (secure password input).

**Parameters:**
- `prompt`: The prompt message to display to the user (e.g., "Enter password: ")

**Return:**
- `string`: The password entered by the user (trimmed of whitespace)
- `error`: Returns error if terminal reading fails

**Behavior:**

1. Displays the prompt to the user
2. Reads password input without echoing characters to the terminal
3. Automatically prints a newline after input
4. Trims and returns the password string

**Security Features:**
- Input is not echoed to the terminal, protecting against shoulder surfing
- Uses system-level terminal control to disable echo
- Input remains hidden while the user types

**Error Handling:**
- Returns error if terminal I/O fails
- May fail if called in environments without proper terminal support

**Example Usage:**

```go
password, err := GetPassword("Enter vault password: ")
if err != nil {
    log.Fatal(err)
}

// Validate password
if password == "" {
    fmt.Println("Password cannot be empty")
}

// Use password
err = Connect(username, password)
if err != nil {
    fmt.Println("Connection failed:", err)
}
```

### Technical Details

- Uses `syscall.Stdin` file descriptor for standard input
- Uses `golang.org/x/term.ReadPassword()` which handles terminal settings
- The `ReadPassword` function doesn't capture the newline character, so one is printed manually
- Input is trimmed of leading/trailing whitespace

### Notes

- This function should be used whenever sensitive input (passwords, API keys, etc.) is required
- The terminal echo state is properly restored after input
- Works cross-platform (Windows, Linux, macOS) through the golang.org/x/term abstraction
