package action

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/urfave/cli"
)

func ComposerGraphAction(g Graph, r GraphRenderer, c *cli.Context) {
	project := c.String("p")
	buildComposerGraph(&g, c, project, 0)

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

func buildComposerGraph(graph *Graph, c *cli.Context, project string, depth int) {
	defer (func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	})()

	if c.Int("d") != -1 && c.Int("d") <= depth {
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

		if addEdge(c, graph, from, to, version) {
			buildComposerGraph(graph, c, root+"/vendor/"+to, depth+1)
		}
	}
}
