package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/urfave/cli"
	"github.com/windler/godepg/action"
	"github.com/windler/godepg/action/composeraction"
	"github.com/windler/godepg/action/configaction"
	"github.com/windler/godepg/action/goaction"
	"github.com/windler/godepg/action/psr4action"
	"github.com/windler/godepg/appcontext"
	"github.com/windler/godepg/dotgraph"
)

func main() {
	app := cli.NewApp()
	app.Author = "Nico Windler"
	app.Copyright = "2017"
	app.Action = func(c *cli.Context) {
		ctx := createContext(c)
		graph := dotgraph.New("godepg")

		file := c.String("file")
		generateGraphFromConfig(file, graph, ctx)
	}
	wd, _ := os.Getwd()
	defaultCfgFile := wd + "/godepg.yml"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file",
			Usage: "the `config-file` to use.",
			Value: defaultCfgFile,
		},
	}
	app.Commands = cli.Commands{
		createComposerCommand(),
		createGOCommand(),
		createPSR4Command(),
	}
	app.Version = "1.0.0"
	app.Description = "Create a dependency graph for your go package."
	app.Usage = "go dependency graph generator"

	app.Run(os.Args)
}

func createContext(c *cli.Context) appcontext.AppContext {
	return appcontext.AppContext{
		Context:      c,
		Bools:        make(map[string]bool),
		Strings:      make(map[string]string),
		Stringslices: make(map[string][]string),
		Ints:         make(map[string]int),
	}
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
				OutputFile: c.String("o"),
			}
			goaction.GenertateGoGraph(graph, renderer, createContext(c))
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "o",
				Usage: "destination `file` to write png to",
			},
			cli.StringFlag{
				Name:  "p",
				Usage: "the `package` to analyze",
			},
			cli.BoolFlag{
				Name:  "n",
				Usage: "hide gos buildin packages",
			},
			cli.IntFlag{
				Name:  "d",
				Value: -1,
				Usage: "limit the depth of the graph",
			},
			cli.StringSliceFlag{
				Name:  "f",
				Usage: "filter package name",
			},
			cli.BoolFlag{
				Name:  "m",
				Usage: "show only subpackages of scanned package",
			},
			cli.StringFlag{
				Name:  "i",
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

func createComposerCommand() cli.Command {
	return cli.Command{
		Name: "php-composer",
		Action: func(c *cli.Context) {
			project := c.String("p")
			if project == "" {
				cli.ShowAppHelpAndExit(c, 2)
			}

			graph := dotgraph.New("php_composer")
			renderer := &dotgraph.PNGRenderer{
				HomeDir:    getDefaultHomeDir(),
				Prefix:     project,
				OutputFile: c.String("o"),
			}
			composeraction.ComposerGraphAction(graph, renderer, createContext(c))
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "o",
				Usage: "destination `file` to write png to",
			},
			cli.StringFlag{
				Name:  "p",
				Usage: "the `project` to analyze",
			},
			cli.StringSliceFlag{
				Name:  "f",
				Usage: "filter project name",
			},
			cli.StringSliceFlag{
				Name:  "s",
				Usage: "dont scan dependencies of package name (pattern)",
			},
			cli.IntFlag{
				Name:  "d",
				Value: -1,
				Usage: "limit the depth of the graph",
			},
		},
	}
}

func createPSR4Command() cli.Command {
	return cli.Command{
		Name: "php-psr4",
		Action: func(c *cli.Context) {
			project := c.String("p")
			if project == "" {
				cli.ShowAppHelpAndExit(c, 2)
			}

			graph := dotgraph.New("php_psr4")
			renderer := &dotgraph.PNGRenderer{
				HomeDir:    getDefaultHomeDir(),
				Prefix:     project,
				OutputFile: c.String("o"),
			}
			psr4action.PSR4GraphAction(graph, renderer, createContext(c))
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "o",
				Usage: "destination `file` to write png to",
			},
			cli.StringFlag{
				Name:  "p",
				Usage: "the `project` to analyze",
			},
			cli.StringSliceFlag{
				Name:  "f",
				Usage: "filter project name",
			},
			cli.StringSliceFlag{
				Name:  "s",
				Usage: "dont scan dependencies of package name (pattern)",
			},
			cli.IntFlag{
				Name:  "d",
				Value: -1,
				Usage: "limit the depth of the graph",
			},
			cli.StringSliceFlag{
				Name:  "e",
				Usage: "exclude folder",
			},
		},
	}
}

func generateGraphFromConfig(file string, g *dotgraph.DotGraph, c action.Context) {
	defer (func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	})()

	fmt.Println("Using: " + file)
	if _, err := os.Stat(file); err != nil {
		panic(err)
	}

	cfg := configaction.CreateContext(file, c)

	for pattern, options := range cfg.Edgestyle {
		g.AddEdgeGraphOptions(pattern, options)
	}

	g.SetNodeGraphOptions(cfg.Nodestyle)
	g.SetGraphOptions(cfg.Graphstyle)

	renderer := &dotgraph.PNGRenderer{
		HomeDir:    getDefaultHomeDir(),
		Prefix:     "godepg",
		OutputFile: c.GetStringFlag("o"),
	}

	switch cfg.Language {
	case "go":
		goaction.GenertateGoGraph(g, renderer, c)
	case "php-composer":
		composeraction.ComposerGraphAction(g, renderer, c)
	case "php-psr4":
		psr4action.PSR4GraphAction(g, renderer, c)
	default:
		panic("No supported languge defined.")
	}
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
