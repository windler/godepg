package composeraction

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/windler/godepg/action"
	"github.com/windler/godepg/action/actionhelper"
	"github.com/windler/godepg/action/matcher"
)

type composerGraphFilter struct {
	from    string
	to      string
	context action.Context
}

var _ action.GraphFilter = &composerGraphFilter{}

func (f composerGraphFilter) GetPreNodeFilters() []action.Matcher {
	return []action.Matcher{
		matcher.NewFilterMatcher(f.from, f.context.GetStringSliceFlag("f")),
		matcher.NewFilterMatcher(f.to, f.context.GetStringSliceFlag("f")),
		matcher.NewFilterMatcher(f.from, f.context.GetStringSliceFlag("s")),
	}
}
func (f composerGraphFilter) GetPostNodeFilters() []action.Matcher {
	return []action.Matcher{}
}

// ComposerGraphAction creates and renders a dependency graph using a prpjects composer.json
func ComposerGraphAction(g action.Graph, r action.GraphRenderer, c action.Context) {
	project := c.GetStringFlag("p")
	buildComposerGraph(g, c, project, 0)

	r.Render(g.String())
}

type composerConfig struct {
	Name    string
	Require map[string]string
}

type composerDependency struct {
	Package string
	Version string
}

func buildComposerGraph(graph action.Graph, c action.Context, project string, depth int) {
	defer (func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	})()

	if c.GetIntFlag("d") != -1 && c.GetIntFlag("d") <= depth {
		return
	}

	composerJSON := &composerConfig{}

	if _, err := os.Stat(project); os.IsNotExist(err) {
		panic("Project not found." + project)
	}

	data, err := ioutil.ReadFile(project + "/composer.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, composerJSON)
	if err != nil {
		panic(err)
	}

	for to, version := range composerJSON.Require {
		from := project
		root := project
		if !strings.Contains(project, "vendor") {
			re := regexp.MustCompile(".*\\/(.+)")
			from = re.FindStringSubmatch(project)[1]
		} else {
			re := regexp.MustCompile("(.*)\\/vendor\\/(.+)")
			from = re.FindStringSubmatch(project)[2]
			root = re.FindStringSubmatch(project)[1]
		}

		filter := &composerGraphFilter{
			from:    from,
			to:      to,
			context: c,
		}

		if actionhelper.AddEdge(graph, from, to, version, filter) {
			buildComposerGraph(graph, c, root+"/vendor/"+to, depth+1)
		}
	}
}
