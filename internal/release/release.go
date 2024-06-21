package release

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func RunRelease() {
	// Step 1: Go mod tidy
	runCommand("go", "mod", "tidy")

	// Step 2: Go build
	runCommand("go", "build", "-o", "kclnr")

	// Step 3: Get version from the built binary
	version := getVersion("./kclnr")

	// Step 4: Check if the tag already exists
	tag := fmt.Sprintf("v%s", version)
	if tagExists(tag) {
		fmt.Printf("Tag %s already exists. Skipping release process.\n", tag)
		return
	}

	// Step 5: Ask for commit message
	commitMessage := getCommitMessage()

	// Step 6: Git add all changes
	runCommand("git", "add", ".")

	// Step 7: Git commit
	runCommand("git", "commit", "-m", commitMessage)

	// Step 8: Git push
	runCommand("git", "push", "origin", "main")

	// Step 9: Git tag
	tagMessage := fmt.Sprintf("Release version %s", version)
	runCommand("git", "tag", "-a", tag, "-m", tagMessage)

	// Step 10: Git push tag
	runCommand("git", "push", "origin", tag)
}

func runCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running command: %s %v\n", name, args)
		os.Exit(1)
	}
}

func getVersion(binaryPath string) string {
	cmd := exec.Command(binaryPath, "version")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting version: %v\n", err)
		os.Exit(1)
	}

	// Assuming the version is the third word in the output
	parts := strings.Fields(out.String())
	if len(parts) < 3 {
		fmt.Fprintln(os.Stderr, "Unexpected version output")
		os.Exit(1)
	}

	return parts[2]
}

func tagExists(tag string) bool {
	cmd := exec.Command("git", "tag", "--list", tag)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking if tag exists: %v\n", err)
		os.Exit(1)
	}
	return strings.TrimSpace(out.String()) == tag
}

func getCommitMessage() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter commit message: ")
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading commit message: %v\n", err)
		os.Exit(1)
	}
	return strings.TrimSpace(message)
}
