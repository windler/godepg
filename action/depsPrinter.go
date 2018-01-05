package action

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/urfave/cli"
)

type depsValues struct {
	Package        string
	Count          int
	Dependencies   []string
	DependencyType string
}

func PrintDeps(depsPkg string, graph *Graph, c *cli.Context) {
	format := c.String("format")
	if format == "" {
		format = "There are {{.Count}} {{.DependencyType}} for package {{.Package}}:\n\n{{range $i, $v := .Dependencies}}{{$i}}: {{$v}}\n{{end}}"
	}

	deps := []string{}
	depsUnquoted := []string{}
	depsValues := &depsValues{
		Package: c.String("info"),
	}

	t := template.New("deps")
	_, err := t.Parse(format)

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	if c.Bool("inverse") {
		deps = (*graph).GetDependents(depsPkg)
		depsValues.DependencyType = "dependents"
	} else {
		deps = (*graph).GetDependencies(depsPkg)
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
}
