package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"reflect"
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
	app.Description = "Create a dependency graph for your go package."
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
			Usage: "shows all packages that depend on the package rather than its dependencies",
		},
		cli.StringFlag{
			Name:  "format",
			Usage: "formats the dependencies output (--info)",
			Value: "There are {{.Count}} {{.DependencyType}} for package {{.Package}}:\n\n{{range $i, $v := .Dependencies}}{{$i}}: {{$v}}\n{{end}}",
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
		return printDeps(c, graph)
	}

	render(graph, dotFile, outFile)
	return nil

}

type depsValues struct {
	Package        string
	Count          int
	Dependencies   []string
	DependencyType string
}

func printDeps(c *cli.Context, graph *graph.Graph) error {
	depsPkg := c.String("info")
	deps := []string{}
	depsUnquoted := []string{}
	depsValues := &depsValues{
		Package: c.String("info"),
	}

	t := template.New("deps")
	_, err := t.Parse(c.String("format"))

	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	if c.Bool("inverse") {
		deps = graph.GetDependents(depsPkg)
		depsValues.DependencyType = "dependents"
	} else {
		deps = graph.GetDependencies(depsPkg)
		depsValues.DependencyType = "dependencies"
	}

	depsValues.Count = len(deps)

	for _, d := range deps {
		depsUnquoted = append(depsUnquoted, strings.Replace(d, "\"", "", -1))
	}
	depsValues.Dependencies = depsUnquoted

	err = t.Execute(os.Stdout, depsValues)
	if err != nil {
		str := reflect.ValueOf(depsValues).Elem()

		fmt.Println(err.Error())
		fmt.Println("")
		fmt.Println("Available fields: ")

		for i := 0; i < str.NumField(); i++ {
			fmt.Println(str.Type().Field(i).Name)
		}
	}
	return nil
}

func render(graph *graph.Graph, dotFile, outFile string) {
	content := graph.GetDotFileContent()
	fmt.Println(content)

	err := ioutil.WriteFile(dotFile, []byte(content), os.ModePerm)
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = exec.Command("dot", "-Tpng", dotFile, "-o", outFile).Output()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = os.Remove(dotFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Written to " + outFile)
}

func createGraph(c *cli.Context, pkg string) *graph.Graph {
	graph := graph.New("godepg")

	data, err := exec.Command("go", "list", "-f", "{{ .ImportPath }}->{{ .Imports }}", pkg+"/...").Output()
	if err != nil {
		log.Fatal(err.Error())
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

	if c.Bool("n") && matcher.NewGoPackagesMatcher(to).Matches() {
		return
	}

	if c.Bool("m") && !matcher.NewSubPackageMatcher(c.String("p"), to).Matches() {
		return
	}

	graph.AddDirectedEdge(from, to)
}
