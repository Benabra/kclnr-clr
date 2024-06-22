package main

import (
	"os"

	"github.com/Benabra/kclnr-clr/cmd"
	"github.com/Benabra/kclnr-clr/internal/release"
	// "github.com/Benabra/kclnr-cli/internal/krew"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "release":
			release.RunRelease()
			return
		case "krew":
			// if len(os.Args) != 6 {
			// 	fmt.Println("Usage: krew <version> <linux-sha256> <darwin-sha256> <windows-sha256>")
			// 	os.Exit(1)
			// }
			// version := os.Args[2]
			// linuxSHA := os.Args[3]
			// darwinSHA := os.Args[4]
			// windowsSHA := os.Args[5]
			// krew.GenerateKrewManifest(version, linuxSHA, darwinSHA, windowsSHA)
			// return
		}
	}

	cmd.Execute()
}
