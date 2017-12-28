package graphviz

import (
	"strings"
)

type Graph struct {
	name  string
	edges []*Edge
	nodes []string
}

type EdgeBuilder struct {
	graph *Graph
	edge  *Edge
}
type Edge struct {
	from string
	to   string
}

func New(name string) *Graph {
	return &Graph{
		name:  name,
		edges: []*Edge{},
		nodes: []string{},
	}
}

func (g *Graph) AddNode(node string) {
	g.nodes = append(g.nodes, node)
}

func (g *Graph) GetDotFileContent() string {
	content := []string{"digraph " + g.name + " {"}

	for _, edge := range g.edges {
		content = append(content, getIdSafeNodeName(edge.from)+"->"+getIdSafeNodeName(edge.to))
	}
	for _, node := range g.nodes {
		content = append(content, getIdSafeNodeName(node))
	}
	content = append(content, "}")

	return strings.Join(content, "\n")
}

func getIdSafeNodeName(id string) string {
	result := strings.TrimSpace(id)
	result = "\"" + result + "\""
	return result
}

func (g *Graph) AddDirectedEdge() *EdgeBuilder {
	return &EdgeBuilder{
		graph: g,
		edge:  &Edge{},
	}
}

func (eb *EdgeBuilder) From(node string) *EdgeBuilder {
	eb.edge.from = node
	return eb
}

func (eb *EdgeBuilder) To(node string) {
	eb.edge.to = node
	eb.graph.edges = append(eb.graph.edges, eb.edge)
}
