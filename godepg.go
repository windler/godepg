package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/urfave/cli"

	"github.com/windler/godepg/config"
	"github.com/windler/godepg/graph"
	"github.com/windler/godepg/http"
	"github.com/windler/godepg/matcher"
)

func main() {

	app := cli.NewApp()
	app.Author = "Nico Windler"
	app.Copyright = "2017"
	app.Action = action
	app.Version = "1.0.0"
	app.Description = "Create a dependency graph for ypur go package."
	app.Usage = "go dependency graph generator"
	app.Flags = []cli.Flag{
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
			Usage: "shows all packages that depend on `package` rather than its dependencies",
		},
	}

	app.Commands = []cli.Command{
		cli.Command{
			Action: func(c *cli.Context) error {
				http.StartWebServer(c.Int64("p"))
				return nil
			},
			Name:  "ws",
			Usage: "starts a webserver to browse all of your generated graphs",
			Flags: []cli.Flag{
				cli.Int64Flag{
					Value: 8000,
					Name:  "p, port",
					Usage: "start webserver on `port`",
				},
			},
		},
	}

	app.Run(os.Args)
}

func action(c *cli.Context) error {
	if c.String("p") == "" {
		cli.ShowAppHelpAndExit(c, 2)
	}

	pkg := c.String("p")
	outFile := c.String("o")
	if c.String("o") == "" {
		pkgName := strings.Replace(pkg, "/", "_", -1)
		pkgName = strings.Replace(pkgName, ".", "_", -1)

		outFile = config.GetDefaultHomeDir() + "/" + pkgName + "_" + time.Now().Format("20060102150405") + ".png"
	}
	dotFile := outFile + ".dot"

	graph := createGraph(c, pkg)

	if c.String("info") != "" {
		printDeps(c, graph)
	} else {
		render(graph, dotFile, outFile)
	}

	return nil
}

func printDeps(c *cli.Context, graph *graph.Graph) {
	depsPkg := c.String("info")
	deps := []string{}

	if c.Bool("inverse") {
		deps = graph.GetDependents(depsPkg)
		fmt.Printf("There are %d dependents for package %s:\n\n", len(deps), depsPkg)
	} else {
		deps = graph.GetDependencies(depsPkg)
		fmt.Printf("There are %d dependencies for package %s:\n\n", len(deps), depsPkg)
	}

	for _, d := range deps {
		fmt.Println(d)
	}
}

func render(graph *graph.Graph, dotFile, outFile string) {
	err := ioutil.WriteFile(dotFile, []byte(graph.GetDotFileContent()), os.ModePerm)
	if err != nil {
		cli.HandleExitCoder(err)
	}

	_, err = exec.Command("dot", "-Tpng", dotFile, "-o", outFile).Output()
	if err != nil {
		cli.HandleExitCoder(err)
	}

	err = os.Remove(dotFile)
	if err != nil {
		cli.HandleExitCoder(err)
	}

	fmt.Println("Written to " + outFile)
}

func createGraph(c *cli.Context, pkg string) *graph.Graph {
	graph := graph.New("godepg")

	data, err := exec.Command("go", "list", "-f", "{{ .ImportPath }}->{{ .Imports }}", pkg+"/...").Output()
	if err != nil {
		cli.HandleExitCoder(err)
	}

	lines := strings.Split(string(data), "]")

	if c.Int("d") >= 0 {
		lines = lines[:c.Int("d")]
	}

	for _, s := range lines {
		packageDeps := strings.Split(s, "->[")
		from := packageDeps[0]
		if len(packageDeps) > 1 {
			for _, to := range strings.Split(packageDeps[1], " ") {
				addEdge(c, graph, from, to)
			}
		}
	}

	return graph
}

func addEdge(c *cli.Context, graph *graph.Graph, from, to string) {
	filterMatcherFrom := matcher.NewFilterMatcher(from, c.StringSlice("f"))
	filterMatcherTo := matcher.NewFilterMatcher(to, c.StringSlice("f"))

	if filterMatcherFrom.Matches() || filterMatcherTo.Matches() {
		return
	}

	graph.AddNode(from)

	if c.Bool("n") && matcher.NewGoPAckagesMatcher(to).Matches() {
		return
	}

	if c.Bool("m") && !matcher.NewSubPackageMatcher(c.String("p"), to).Matches() {
		return
	}

	graph.AddDirectedEdge(from, to)
}
