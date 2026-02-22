package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	fmt.Println("🕊️  Salam Installer")
	fmt.Println("==================")

	binaryName := "salam"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	// 1. Find the binary to install
	// We expect the 'salam' binary to be in the same folder as the installer
	sourcePath, _ := os.Executable()
	sourceDir := filepath.Dir(sourcePath)
	sourceBinary := filepath.Join(sourceDir, binaryName)

	if _, err := os.Stat(sourceBinary); os.IsNotExist(err) {
		// Fallback: search in build directory
		sourceBinary = filepath.Join(sourceDir, "build", binaryName)
		if _, err := os.Stat(sourceBinary); os.IsNotExist(err) {
			fmt.Printf("❌ Could not find %s to install.\n", binaryName)
			fmt.Println("Please make sure the binary is in the same directory as the installer.")
			os.Exit(1)
		}
	}

	// 2. Determine installation directory
	var installDir string
	if runtime.GOOS == "windows" {
		home, _ := os.UserHomeDir()
		installDir = filepath.Join(home, ".salam", "bin")
	} else {
		installDir = "/usr/local/bin"
	}

	// 3. Create directory
	if err := os.MkdirAll(installDir, 0755); err != nil {
		fmt.Printf("❌ Failed to create installation directory: %v\n", err)
		os.Exit(1)
	}

	// 4. Copy binary
	targetPath := filepath.Join(installDir, binaryName)
	input, err := os.ReadFile(sourceBinary)
	if err != nil {
		fmt.Printf("❌ Failed to read source binary: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(targetPath, input, 0755); err != nil {
		if runtime.GOOS != "windows" && os.IsPermission(err) {
			fmt.Println("🔐 Permission denied. Retrying with sudo...")
			cmd := exec.Command("sudo", "cp", sourceBinary, targetPath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Printf("❌ Failed to install with sudo: %v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Printf("❌ Failed to write binary to %s: %v\n", targetPath, err)
			os.Exit(1)
		}
	}

	// 5. Windows specific: Update PATH
	if runtime.GOOS == "windows" {
		fmt.Println("📝 Adding Salam to your PATH...")
		// Use powershell to update user path persistently
		psCmd := fmt.Sprintf("[Environment]::SetEnvironmentVariable('Path', [Environment]::GetEnvironmentVariable('Path', 'User') + ';%s', 'User')", installDir)
		cmd := exec.Command("powershell", "-Command", psCmd)
		if err := cmd.Run(); err != nil {
			fmt.Printf("⚠️  Failed to update PATH automatically: %v\n", err)
			fmt.Printf("Please add %s to your PATH manually.\n", installDir)
		} else {
			fmt.Println("✅ PATH updated successfully!")
		}
	}

	fmt.Println("\n🎉 Installation successful!")
	fmt.Printf("Salam has been installed to: %s\n", targetPath)
	fmt.Println("\nPlease restart your terminal and type 'salam' to get started!")
}
