package dotgraph

import (
	"strings"
)

//Graph represents the dependency graph of a package
type dotGraph struct {
	name  string
	edges map[string][]string
	nodes []string
}

//New creates a new Graph with a given name
func New(name string) *dotGraph {
	return &dotGraph{
		name:  name,
		edges: make(map[string][]string),
		nodes: []string{},
	}
}

//AddNode add a node to string. There does not have to be an edge for a node.
func (g dotGraph) AddNode(node string) {
	contains := false
	new := getIDSafeNodeName(node)

	for _, n := range g.nodes {
		if n == new {
			contains = true
			break
		}
	}
	if !contains {
		g.nodes = append(g.nodes, new)
	}
}

//GetDotFileContent create the content of a dot-file (graphviz)
func (g dotGraph) String() string {
	content := []string{"digraph " + g.name + " {"}

	for from, deps := range g.edges {
		for _, to := range deps {
			if from != `""` && to != `""` {
				content = append(content, from+"->"+to)
			}
		}
	}
	for _, node := range g.nodes {
		content = append(content, node)
	}
	content = append(content, "}")

	return strings.Join(content, "\n")
}

func getIDSafeNodeName(id string) string {
	result := strings.TrimSpace(id)
	if strings.HasSuffix(result, "/") {
		result = result[0 : len(result)-1]
	}
	result = "\"" + result + "\""

	return result
}

//AddDirectedEdge adds an directed edge for two nodes to the graph
func (g dotGraph) AddDirectedEdge(from, to string) {
	saveFrom := getIDSafeNodeName(from)
	saveTo := getIDSafeNodeName(to)
	if _, found := g.edges[saveFrom]; !found {
		g.edges[saveFrom] = []string{}
	}
	g.edges[saveFrom] = append(g.edges[saveFrom], saveTo)
}

//GetDependencies returns alls direct dipendencies for a package within the graph
func (g dotGraph) GetDependencies(pkg string) []string {
	dependencies := []string{}

	for from, deps := range g.edges {
		if from == getIDSafeNodeName(pkg) {
			dependencies = deps
		}
	}
	return dependencies
}

//GetDependents returns all packages that directly depend on the given package within the graph
func (g dotGraph) GetDependents(pkg string) []string {
	dependents := []string{}
loop:
	for from, deps := range g.edges {
		for _, to := range deps {
			if to == getIDSafeNodeName(pkg) {
				dependents = append(dependents, from)
				continue loop
			}
		}
	}
	return dependents
}
