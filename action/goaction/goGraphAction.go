package goaction

import (
	"log"
	"os/exec"
	"strings"

	"github.com/windler/godepg/action"
	"github.com/windler/godepg/action/actionhelper"
	"github.com/windler/godepg/action/matcher"
)

type goGraphFilter struct {
	from    string
	to      string
	context action.Context
}

var _ action.GraphFilter = &goGraphFilter{}

func (f goGraphFilter) GetPreNodeFilters() []action.Matcher {
	return []action.Matcher{
		matcher.NewFilterMatcher(f.from, f.context.GetStringSliceFlag("f")),
		matcher.NewFilterMatcher(f.to, f.context.GetStringSliceFlag("f")),
	}
}
func (f goGraphFilter) GetPostNodeFilters() []action.Matcher {
	res := []action.Matcher{}
	if f.context.GetBoolFlag("n") {
		res = append(res, matcher.NewGoPackagesMatcher(f.to))
	}

	if f.context.GetBoolFlag("m") {
		res = append(res, matcher.NewSubPackageMatcher(f.context.GetStringFlag("p"), f.to))
	}

	return res
}

//GenertateGoGraph generates a dependency graph for a go package
func GenertateGoGraph(g action.Graph, r action.GraphRenderer, c action.Context) {
	pkg := c.GetStringFlag("p")
	buildGoGraph(g, c, pkg)

	if c.GetStringFlag("info") != "" {
		printDeps(c.GetStringFlag("info"), &g, c)
	} else {
		r.Render(g.String())
	}
}

func buildGoGraph(graph action.Graph, c action.Context, pkg string) {
	data, err := exec.Command("go", "list", "-f", "{{ .ImportPath }}->{{ .Imports }}", pkg+"/...").Output()
	if err != nil {
		log.Fatal(err.Error())
	}

	lines := strings.Split(string(data), "]")

	if c.GetIntFlag("d") >= 0 {
		lines = lines[:c.GetIntFlag("d")]
	}

	for _, s := range lines {
		packageDeps := strings.Split(s, "->[")
		from := packageDeps[0]
		if len(packageDeps) > 1 {
			for _, to := range strings.Split(packageDeps[1], " ") {
				filter := &goGraphFilter{
					from:    from,
					to:      to,
					context: c,
				}
				actionhelper.AddEdge(graph, from, to, "", filter)
			}
		}
	}
}
