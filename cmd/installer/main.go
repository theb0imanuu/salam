package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

//go:embed salam.exe
var winBinary []byte

//go:embed salam
var unixBinary []byte

func main() {
	fmt.Println("🕊️  Salam Installer")
	fmt.Println("==================")

	binaryName := "salam"
	var binaryData []byte

	if runtime.GOOS == "windows" {
		binaryName += ".exe"
		binaryData = winBinary
	} else {
		binaryData = unixBinary
	}

	if len(binaryData) == 0 {
		fmt.Printf("❌ Error: Embedded binary for %s is empty.\n", runtime.GOOS)
		fmt.Println("This installer was built incorrectly.")
		pause()
		os.Exit(1)
	}

	// 1. Determine installation directory
	var installDir string
	if runtime.GOOS == "windows" {
		home, _ := os.UserHomeDir()
		installDir = filepath.Join(home, ".salam", "bin")
	} else {
		installDir = "/usr/local/bin"
	}

	// 2. Create directory
	if err := os.MkdirAll(installDir, 0755); err != nil {
		fmt.Printf("❌ Failed to create installation directory: %v\n", err)
		pause()
		os.Exit(1)
	}

	// 3. Write binary
	targetPath := filepath.Join(installDir, binaryName)
	if err := os.WriteFile(targetPath, binaryData, 0755); err != nil {
		if runtime.GOOS != "windows" && os.IsPermission(err) {
			fmt.Println("🔐 Permission denied. Please run with sudo or check permissions.")
		}
		fmt.Printf("❌ Failed to write binary to %s: %v\n", targetPath, err)
		pause()
		os.Exit(1)
	}

	// 4. Windows specific: Update PATH
	if runtime.GOOS == "windows" {
		fmt.Println("📝 Adding Salam to your PATH...")
		if err := updateWindowsPath(installDir); err != nil {
			fmt.Printf("⚠️  Failed to update PATH automatically: %v\n", err)
			fmt.Printf("Please add %s to your PATH manually.\n", installDir)
		} else {
			fmt.Println("✅ PATH updated successfully!")
		}
	}

	fmt.Println("\n🎉 Installation successful!")
	fmt.Printf("Salam has been installed to: %s\n", targetPath)

	if runtime.GOOS == "windows" {
		fmt.Println("\nIMPORTANT: Please close this terminal and open a NEW one for the changes to take effect.")
		fmt.Println("Then just type 'salam' to get started!")
	} else {
		fmt.Println("\nType 'salam' to get started!")
	}

	pause()
}

func updateWindowsPath(newPath string) error {
	// 1. Get current path
	getCmd := exec.Command("powershell", "-Command", "[Environment]::GetEnvironmentVariable('Path', 'User')")
	out, err := getCmd.Output()
	if err != nil {
		return err
	}
	currentPath := strings.TrimSpace(string(out))

	// 2. Check if already exists
	paths := strings.Split(currentPath, ";")
	for _, p := range paths {
		if strings.EqualFold(p, newPath) {
			return nil // Already in PATH
		}
	}

	// 3. Append and set
	updatedPath := currentPath
	if updatedPath != "" && !strings.HasSuffix(updatedPath, ";") {
		updatedPath += ";"
	}
	updatedPath += newPath

	setCmd := exec.Command("powershell", "-Command", fmt.Sprintf("[Environment]::SetEnvironmentVariable('Path', '%s', 'User')", updatedPath))
	return setCmd.Run()
}

func pause() {
	if runtime.GOOS == "windows" {
		fmt.Println("\nPress Enter to exit...")
		fmt.Scanln()
	}
}
