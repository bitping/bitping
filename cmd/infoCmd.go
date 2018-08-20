package cmd

import (
	"fmt"
	"runtime"

	"github.com/codegangsta/cli"
)

var (
	// AppName of the app
	AppName = "unset"
	// Version of the app
	Version = "unset"
	// BuildTime when built
	BuildTime = "unset"
	// Commit when built
	Commit string
	// Arch running
	Arch = runtime.GOARCH
	// OS of the binary
	OS = runtime.GOOS
	// Branch of build
	Branch string

	// GoVersion is the go version this was compiled against
	GoVersion = runtime.Version()
)

// InfoCmd provides info about the build
var InfoCmd = cli.Command{
	Name:   "info",
	Usage:  "information about bitping",
	Action: displayBitpingInfo,
}

func displayBitpingInfo(c *cli.Context) {
	msg := fmt.Sprintf(`%s
OS: %s
GoVersion: %v
Commit: %s
Branch: %s
Built at %s
`, AppName, OS, GoVersion, BuildTime, Branch, Commit)

	fmt.Println(msg)
}
