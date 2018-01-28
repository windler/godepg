package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"regexp"

	"github.com/urfave/cli"
	"github.com/windler/godepg/action"
	"github.com/windler/godepg/action/composeraction"
	"github.com/windler/godepg/action/goaction"
	"github.com/windler/godepg/dotgraph"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	app := cli.NewApp()
	app.Author = "Nico Windler"
	app.Copyright = "2017"
	app.Action = func(c *cli.Context) {
		ctx := createContext(c)
		graph := dotgraph.New("godepg")

		GenerateGraphFromConfig(c.String("file"), graph, ctx)
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file",
			Usage: "the `config-file` to use.",
			Value: "godepg.yml",
		},
	}
	app.Commands = cli.Commands{
		createComposerCommand(),
		createGOCommand(),
	}
	app.Version = "1.0.0"
	app.Description = "Create a dependency graph for your go package."
	app.Usage = "go dependency graph generator"

	app.Run(os.Args)
}

func createContext(c *cli.Context) AppContext {
	return AppContext{
		context:      c,
		bools:        make(map[string]bool),
		strings:      make(map[string]string),
		stringslices: make(map[string][]string),
		ints:         make(map[string]int),
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

type config struct {
	Language   string
	Project    string
	Filter     []string
	Depth      int
	StopAt     []string
	Output     string
	Edgestyle  map[string]dotgraph.DotGraphOptions
	Nodestyle  dotgraph.DotGraphOptions
	Graphstyle dotgraph.DotGraphOptions
}

//GenerateGraphFromConfig reads a config file and generates a graph based on the config
func GenerateGraphFromConfig(file string, g *dotgraph.DotGraph, c action.Context) {
	defer (func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	})()

	if _, err := os.Stat(file); err != nil {
		panic(err)
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	cfg := &config{}

	yaml.Unmarshal(data, cfg)

	prepareContext(file, cfg, c)

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
	default:
		panic("No supported languge defined.")
	}
}

func prepareContext(file string, cfg *config, context action.Context) {
	context.SetStringSliceFlag("f", cfg.Filter)
	context.SetStringSliceFlag("s", cfg.StopAt)

	context.SetStringFlag("p", cfg.Project)
	if cfg.Project == "" {
		re := regexp.MustCompile("(.*)\\/(.+)")
		projectRoot := re.FindStringSubmatch(file)[1]
		context.SetStringFlag("p", projectRoot)
	}

	context.SetIntFlag("d", cfg.Depth)
	if cfg.Depth == 0 {
		context.SetIntFlag("d", -1)
	}

	if cfg.Output != "" {
		context.SetStringFlag("o", cfg.Output)
	}
}

//AppContext provides app flags
type AppContext struct {
	context      *cli.Context
	strings      map[string]string
	ints         map[string]int
	stringslices map[string][]string
	bools        map[string]bool
}

//GetStringFlag gets the value of a string flag
func (ac AppContext) GetStringFlag(flag string) string {
	if res, found := ac.strings[flag]; found {
		return res
	}
	return ac.context.String(flag)
}

//GetStringSliceFlag gets all values for a slice flag
func (ac AppContext) GetStringSliceFlag(flag string) []string {
	if res, found := ac.stringslices[flag]; found {
		return res
	}
	return ac.context.StringSlice(flag)
}

//GetIntFlag gets an int-value for a flag
func (ac AppContext) GetIntFlag(flag string) int {
	if res, found := ac.ints[flag]; found {
		return res
	}
	return ac.context.Int(flag)
}

//GetBoolFlag gets a bool-value for a flag
func (ac AppContext) GetBoolFlag(flag string) bool {
	if res, found := ac.bools[flag]; found {
		return res
	}
	return ac.context.Bool(flag)
}

//SetStringFlag sets a string flag
func (ac AppContext) SetStringFlag(flag, value string) {
	ac.strings[flag] = value
}

//SetStringSliceFlag sets a stringslice flag
func (ac AppContext) SetStringSliceFlag(flag string, value []string) {
	ac.stringslices[flag] = value
}

//SetIntFlag sets a int flag
func (ac AppContext) SetIntFlag(flag string, value int) {
	ac.ints[flag] = value
}

//SetBoolFlag sets a bool flag
func (ac AppContext) SetBoolFlag(flag string, value bool) {
	ac.bools[flag] = value
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
