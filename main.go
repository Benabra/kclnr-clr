package main

import (
	"os"

	"github.com/Benabra/kclnr-clr/cmd"
	"github.com/Benabra/kclnr-clr/internal/release"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "release" {
		release.RunRelease()
		return
	}
	cmd.Execute()
}
