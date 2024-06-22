package krew

import (
	"fmt"
	"os"
	"text/template"
)

const krewTemplate = `apiVersion: krew.googlecontainertools.github.com/v1alpha2
kind: Plugin
metadata:
  name: kclnr
spec:
  version: "{{.Version}}"
  shortDescription: "Manage Kubernetes contexts and nodes easily"
  description: |
    kclnr is a command-line tool to simplify the management of Kubernetes contexts and nodes.
  platforms:
    - selector:
        matchLabels:
          os: linux
          arch: amd64
      uri: "https://github.com/Benabra/kclnr-cli/releases/download/v{{.Version}}/kclnr-linux-amd64.tar.gz"
      sha256: "{{.LinuxSHA256}}"
      bin: kclnr
      files:
        - from: "*"
          to: "."
    - selector:
        matchLabels:
          os: darwin
          arch: amd64
      uri: "https://github.com/Benabra/kclnr-cli/releases/download/v{{.Version}}/kclnr-darwin-amd64.tar.gz"
      sha256: "{{.DarwinSHA256}}"
      bin: kclnr
      files:
        - from: "*"
          to: "."
    - selector:
        matchLabels:
          os: windows
          arch: amd64
      uri: "https://github.com/Benabra/kclnr-cli/releases/download/v{{.Version}}/kclnr-windows-amd64.zip"
      sha256: "{{.WindowsSHA256}}"
      bin: kclnr.exe
      files:
        - from: "*"
          to: "."
`

type ManifestData struct {
	Version       string
	LinuxSHA256   string
	DarwinSHA256  string
	WindowsSHA256 string
}

func GenerateKrewManifest(version, linuxSHA256, darwinSHA256, windowsSHA256 string) {
	data := ManifestData{
		Version:       version,
		LinuxSHA256:   linuxSHA256,
		DarwinSHA256:  darwinSHA256,
		WindowsSHA256: windowsSHA256,
	}

	tmpl, err := template.New("krew").Parse(krewTemplate)
	if err != nil {
		fmt.Println("Error creating template:", err)
		return
	}

	err = os.MkdirAll("krew", 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	file, err := os.Create("krew/kclnr.yaml")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	err = tmpl.Execute(file, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}

	fmt.Println("Krew plugin manifest generated at krew/kclnr.yaml")
}
