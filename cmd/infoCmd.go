package cmd

import (
	"fmt"
	"os"
	"runtime"
	"text/tabwriter"

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

	w := tabwriter.NewWriter(os.Stdout, 10, 4, 4, ' ', tabwriter.DiscardEmptyColumns)

	fmt.Fprintf(w, "Binary\t%s\n", au.Green(AppName))
	fmt.Fprintf(w, "OS\t%s\n", au.Green(OS))
	fmt.Fprintf(w, "Branch\t%s\n", au.Green(Branch))
	fmt.Fprintf(w, "Commit\t%s\n", au.Green(Commit))
	fmt.Fprintf(w, "Go version\t%s\n", au.Green(GoVersion))
	fmt.Fprintf(w, "Build time\t%s\n", au.Green(BuildTime))

	w.Flush()
	// 	msg := fmt.Sprintf(`%s
	// OS: %s
	// Go version: %v
	// Commit: %s
	// Branch: %s
	// Built at %s
	// `,
	// 		au.Green(AppName),
	// 		au.Green(OS),
	// 		au.Green(GoVersion),
	// 		au.Green(Commit),
	// 		au.Green(Branch),
	// 		au.Green(BuildTime))

	// 	fmt.Println(msg)
}
