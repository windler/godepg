package main

import (
	"log"
	"os"
	"os/user"

	"github.com/urfave/cli"
	"github.com/windler/godepg/action/composeraction"
	"github.com/windler/godepg/action/goaction"
	"github.com/windler/godepg/dotgraph"
)

func main() {
	app := cli.NewApp()
	app.Author = "Nico Windler"
	app.Copyright = "2017"
	app.Action = func(c *cli.Context) {
		pkg := c.String("p")
		if pkg == "" {
			cli.ShowAppHelpAndExit(c, 2)
		}

		graph := dotgraph.New("godepg")
		renderer := &dotgraph.PNGRenderer{
			HomeDir:    getDefaultHomeDir(),
			Prefix:     pkg,
			OutputFile: c.String("output"),
		}
		goaction.GenertateGoGraph(graph, renderer, &AppContext{context: c})
	}
	app.Commands = cli.Commands{
		createPHPCommand(),
		createGOCommand(),
	}
	app.Version = "1.0.0"
	app.Description = "Create a dependency graph for your go package."
	app.Usage = "go dependency graph generator"

	app.Run(os.Args)
}

func createGOCommand() cli.Command {
	return cli.Command{
		Name: "go",
		Action: func(c *cli.Context) {
			pkg := c.String("p")
			if pkg == "" {
				cli.ShowAppHelpAndExit(c, 2)
			}

			graph := dotgraph.New("godepg")
			renderer := &dotgraph.PNGRenderer{
				HomeDir:    getDefaultHomeDir(),
				Prefix:     pkg,
				OutputFile: c.String("output"),
			}
			goaction.GenertateGoGraph(graph, renderer, &AppContext{context: c})
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "o, output",
				Usage: "destination `file` to write png to",
			},
			cli.StringFlag{
				Name:  "p, package",
				Usage: "the `package` to analyze",
			},
			cli.BoolFlag{
				Name:  "n, no-go-packages",
				Usage: "hide gos buildin packages",
			},
			cli.IntFlag{
				Name:  "d, depth",
				Value: -1,
				Usage: "limit the depth of the graph",
			},
			cli.StringSliceFlag{
				Name:  "f, filter",
				Usage: "filter package name",
			},
			cli.BoolFlag{
				Name:  "m, my-packages-only",
				Usage: "show only subpackages of scanned package",
			},
			cli.StringFlag{
				Name:  "i, info",
				Usage: "shows the dependencies for a `package`",
			},
			cli.BoolFlag{
				Name:  "inverse",
				Usage: "shows all packages that depend on the package rather than its dependencies",
			},
			cli.StringFlag{
				Name:  "format",
				Usage: "formats the dependencies output (--info)",
			},
		},
	}
}

func createPHPCommand() cli.Command {
	return cli.Command{
		Name: "php",
		Action: func(c *cli.Context) {
			project := c.String("p")
			if project == "" {
				cli.ShowAppHelpAndExit(c, 2)
			}

			graph := dotgraph.New("php_composer")
			renderer := &dotgraph.PNGRenderer{
				HomeDir:    getDefaultHomeDir(),
				Prefix:     project,
				OutputFile: c.String("output"),
			}
			composeraction.ComposerGraphAction(graph, renderer, &AppContext{context: c})
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "o, output",
				Usage: "destination `file` to write png to",
			},
			cli.StringFlag{
				Name:  "p, project",
				Usage: "the `project` to analyze",
			},
			cli.StringSliceFlag{
				Name:  "f, filter",
				Usage: "filter project name",
			},
			cli.StringSliceFlag{
				Name:  "s, stop-at",
				Usage: "dont scan dependencies of package name (pattern)",
			},
			cli.IntFlag{
				Name:  "d, depth",
				Value: -1,
				Usage: "limit the depth of the graph",
			},
		},
	}
}

//AppContext provides app flags
type AppContext struct {
	context *cli.Context
}

//GetStringFlag gets the value of a string flag
func (ac AppContext) GetStringFlag(flag string) string {
	return ac.context.String(flag)
}

//GetStringSliceFlag gets all values for a slice flag
func (ac AppContext) GetStringSliceFlag(flag string) []string {
	return ac.context.StringSlice(flag)
}

//GetIntFlag gets an int-value for a flag
func (ac AppContext) GetIntFlag(flag string) int {
	return ac.context.Int(flag)
}

//GetBoolFlag gets a bool-value for a flag
func (ac AppContext) GetBoolFlag(flag string) bool {
	return ac.context.Bool(flag)
}

func getDefaultHomeDir() string {
	usr, _ := user.Current()
	home := usr.HomeDir + "/" + "godepg"

	if _, err := os.Stat(home); os.IsNotExist(err) {
		e := os.Mkdir(home, os.ModePerm)
		if e != nil {
			log.Fatal("Cannot create folder ", err)
		}
	}

	return home
}
