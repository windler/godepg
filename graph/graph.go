package graph

import (
	"strings"
)

type Graph struct {
	name  string
	edges map[string][]string
	nodes []string
}

func New(name string) *Graph {
	return &Graph{
		name:  name,
		edges: make(map[string][]string),
		nodes: []string{},
	}
}

func (g *Graph) AddNode(node string) {
	g.nodes = append(g.nodes, getIdSafeNodeName(node))
}

func (g *Graph) GetDotFileContent() string {
	content := []string{"digraph " + g.name + " {"}

	for from, deps := range g.edges {
		for _, to := range deps {
			content = append(content, from+"->"+to)
		}
	}
	for _, node := range g.nodes {
		content = append(content, node)
	}
	content = append(content, "}")

	return strings.Join(content, "\n")
}

func getIdSafeNodeName(id string) string {
	result := strings.TrimSpace(id)
	if strings.HasSuffix(result, "/") {
		result = result[0 : len(result)-1]
	}
	result = "\"" + result + "\""

	return result
}

func (g *Graph) AddDirectedEdge(from, to string) {
	saveFrom := getIdSafeNodeName(from)
	saveTo := getIdSafeNodeName(to)
	if _, found := g.edges[saveFrom]; !found {
		g.edges[saveFrom] = []string{}
	}
	g.edges[saveFrom] = append(g.edges[saveFrom], saveTo)
}

func (g *Graph) GetDependencies(pkg string) []string {
	dependencies := []string{}

	for from, deps := range g.edges {
		if from == getIdSafeNodeName(pkg) {
			dependencies = deps
		}
	}
	return dependencies
}

func (g *Graph) GetDependents(pkg string) []string {
	dependents := []string{}
loop:
	for from, deps := range g.edges {
		for _, to := range deps {
			if to == getIdSafeNodeName(pkg) {
				dependents = append(dependents, from)
				continue loop
			}
		}
	}
	return dependents
}
