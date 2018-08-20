package cmd

import (
	"fmt"

	"github.com/auser/bitping/database"
	"github.com/codegangsta/cli"
)

var sharedDbFlags = append([]cli.Flag{}, cli.StringFlag{
	Name:  "db-path,p",
	Usage: "database path",
	Value: database.GetDatabasePath(),
})

var queryCmd = cli.Command{
	Name:  "query",
	Usage: "query the databae",
	Action: func(c *cli.Context) {
		fmt.Printf("Querying the database")
	},
}

// DbCmd is the main database command
var DbCmd = cli.Command{
	Name:  "db",
	Usage: "Database actions",
	Flags: sharedDbFlags,
	Subcommands: []cli.Command{
		queryCmd,
	},
}
