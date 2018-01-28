package psr4action

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/windler/godepg/action"
	"github.com/windler/godepg/action/actionhelper"
	"github.com/windler/godepg/action/matcher"
)

type psr4GraphFilter struct {
	from    string
	to      string
	context action.Context
}

var _ action.GraphFilter = &psr4GraphFilter{}

func (f psr4GraphFilter) GetPreNodeFilters() []action.Matcher {
	return []action.Matcher{
		matcher.NewFilterMatcher(f.from, f.context.GetStringSliceFlag("f")),
		matcher.NewFilterMatcher(f.to, f.context.GetStringSliceFlag("f")),
		matcher.NewFilterMatcher(f.from, f.context.GetStringSliceFlag("s")),
	}
}
func (f psr4GraphFilter) GetPostNodeFilters() []action.Matcher {
	return []action.Matcher{}
}

// PSR4GraphAction creates and renders a dependency graph using psr4 imports
func PSR4GraphAction(g action.Graph, r action.GraphRenderer, c action.Context) {
	project := c.GetStringFlag("p")
	buildPSR4Graph(g, c, project, 0)

	r.Render(g.String())
}

func buildPSR4Graph(graph action.Graph, c action.Context, project string, depth int) {
	defer (func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	})()

	if c.GetIntFlag("d") != -1 && c.GetIntFlag("d") <= depth {
		return
	}

	if _, err := os.Stat(project); os.IsNotExist(err) {
		panic("Project not found." + project)
	}

	deps := make(map[string][]string)
	filepath.Walk(project, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".php" {
			extractDeps(path, deps)
		}
		return nil
	})

	for from, toSlice := range deps {
		//add node anyways?
		for _, to := range toSlice {
			filter := &psr4GraphFilter{
				from:    from,
				to:      to,
				context: c,
			}
			actionhelper.AddEdge(graph, from, to, "", filter)
		}
	}

}

func extractDeps(path string, deps map[string][]string) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(f)

	ns := ""
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err.Error() == "EOF" {
				return
			}
			panic(err)
		}

		re := regexp.MustCompile(".*namespace (.+).*;")
		if match := re.FindStringSubmatch(string(line)); len(match) > 0 {
			if _, found := deps[match[1]]; !found {
				ns = strings.Replace(match[1], "\\", "/", -1)
				deps[ns] = []string{}
			}
		}

		re = regexp.MustCompile(".*use (.+)\\sas")
		if match := re.FindStringSubmatch(string(line)); len(match) > 0 {
			deps[ns] = append(deps[ns], extractPackageFromUseClause(strings.Replace(match[1], "\\", "/", -1)))
		} else {
			re = regexp.MustCompile(".*use (.+);")
			if match := re.FindStringSubmatch(string(line)); len(match) > 0 {
				deps[ns] = append(deps[ns], extractPackageFromUseClause(strings.Replace(match[1], "\\", "/", -1)))
			}
		}

		if re := regexp.MustCompile("(.*)class(.+)"); len(re.FindStringSubmatch(string(line))) != 0 {
			return
		}
	}
}

func extractPackageFromUseClause(use string) string {
	re := regexp.MustCompile("(.*)\\/.+")
	if match := re.FindStringSubmatch(use); len(match) > 0 {
		return match[1]
	}

	return use
}
