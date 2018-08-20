package cmd

import (
	"fmt"
	"runtime"

	"github.com/codegangsta/cli"
	"github.com/logrusorgru/aurora"
)

var (
	au aurora.Aurora
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
	au = aurora.NewAurora(!c.GlobalBool("nocolor"))
	msg := fmt.Sprintf(`%s
OS: %s
Go version: %v
Commit: %s
Branch: %s
Built at %s
`,
		au.Green(AppName),
		au.Green(OS),
		au.Green(GoVersion),
		au.Green(Commit),
		au.Green(Branch),
		au.Green(BuildTime))

	fmt.Println(msg)
}
