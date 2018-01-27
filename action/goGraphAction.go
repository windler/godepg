package action

import (
	"log"
	"os/exec"
	"strings"

	"github.com/urfave/cli"
	"github.com/windler/godepg/action/matcher"
)

func GoGraphAction(g Graph, r GraphRenderer, c *cli.Context) {
	pkg := c.String("p")
	buildGraph(&g, c, pkg)

	if c.String("info") != "" {
		PrintDeps(c.String("info"), &g, c)
	} else {
		r.Render(g.String())
	}
}

func buildGraph(graph *Graph, c *cli.Context, pkg string) {
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
				addEdge(c, graph, from, to, "")
			}
		}
	}
}

func addEdge(c *cli.Context, graph *Graph, from, to, description string) bool {
	filterMatcherFrom := matcher.NewFilterMatcher(from, c.StringSlice("f"))
	filterMatcherTo := matcher.NewFilterMatcher(to, c.StringSlice("f"))
	filterMatcherStop := matcher.NewFilterMatcher(from, c.StringSlice("s"))

	if filterMatcherFrom.Matches() || filterMatcherTo.Matches() || filterMatcherStop.Matches() {
		return false
	}

	(*graph).AddNode(from)

	if c.Bool("n") && matcher.NewGoPackagesMatcher(to).Matches() {
		return false
	}

	if c.Bool("m") && !matcher.NewSubPackageMatcher(c.String("p"), to).Matches() {
		return false
	}

	(*graph).AddDirectedEdge(from, to, description)
	return true
}
