package release

import (
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

	// Step 4: Git add all changes
	runCommand("git", "add", ".")

	// Step 5: Git commit
	commitMessage := "Fix module import paths and update project structure"
	runCommand("git", "commit", "-m", commitMessage)

	// Step 6: Git push
	runCommand("git", "push", "origin", "main")

	// Step 7: Git tag
	tag := fmt.Sprintf("v%s", version)
	tagMessage := fmt.Sprintf("Release version %s", version)
	runCommand("git", "tag", "-a", tag, "-m", tagMessage)

	// Step 8: Git push tag
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