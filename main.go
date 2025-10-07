package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: wowcli <command> [options]")
		fmt.Println("Available commands:")
		fmt.Println("  char-info - Display character and account information")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "char-info":
		charInfo()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Available commands:")
		fmt.Println("  char-info - Display character and account information")
		os.Exit(1)
	}
}

func charInfo() {
	// Parse flags for char-info command
	charInfoCmd := flag.NewFlagSet("char-info", flag.ExitOnError)
	customPath := charInfoCmd.String("path", "", "Custom path to WoW Account directory")
	charInfoCmd.Parse(os.Args[2:])

	// Default WoW installation path
	basePath := `C:\Program Files (x86)\World of Warcraft\_retail_\WTF\Account\`

	// Use custom path if provided
	if *customPath != "" {
		basePath = *customPath
		if basePath[len(basePath)-1] != '\\' && basePath[len(basePath)-1] != '/' {
			basePath += string(os.PathSeparator)
		}
	}

	// Check if the directory exists
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		fmt.Printf("WoW account directory not found at: %s\n", basePath)
		fmt.Println("Please ensure World of Warcraft is installed at the default location.")
		os.Exit(1)
	}

	// Read the Account directory
	accounts, err := os.ReadDir(basePath)
	if err != nil {
		fmt.Printf("Error reading account directory: %v\n", err)
		os.Exit(1)
	}

	if len(accounts) == 0 {
		fmt.Println("No accounts found.")
		return
	}

	fmt.Println("World of Warcraft Characters and Accounts")
	fmt.Println("==========================================")
	fmt.Println()

	// Iterate through accounts
	for _, account := range accounts {
		if !account.IsDir() {
			continue
		}

		accountName := account.Name()
		fmt.Printf("Account: %s\n", accountName)

		// Look for realm/server folders in SavedVariables or iterate through the account structure
		accountPath := basePath + accountName + string(os.PathSeparator)

		// Check for realm folders (typically in the account directory)
		servers, err := os.ReadDir(accountPath)
		if err != nil {
			fmt.Printf("  Error reading account directory: %v\n", err)
			continue
		}

		characterCount := 0
		for _, server := range servers {
			if !server.IsDir() {
				continue
			}

			serverName := server.Name()
			// Skip SavedVariables and other non-server directories
			if serverName == "SavedVariables" || serverName == "macros-cache.txt" {
				continue
			}

			serverPath := accountPath + serverName + string(os.PathSeparator)
			characters, err := os.ReadDir(serverPath)
			if err != nil {
				continue
			}

			for _, char := range characters {
				if !char.IsDir() {
					continue
				}

				characterName := char.Name()
				fmt.Printf("  %s - %s\n", serverName, characterName)
				characterCount++
			}
		}

		if characterCount == 0 {
			fmt.Println("  No characters found")
		}
		fmt.Println()
	}
}
