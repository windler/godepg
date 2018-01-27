package dotgraph

import (
	"regexp"
	"strings"
)

//DotGraph create a directed graph in dot notation
type DotGraph struct {
	name         string
	edges        map[string][]edge
	edgeOptions  map[string]DotGraphOptions
	nodeOptions  DotGraphOptions
	graphOptions DotGraphOptions
}

type edge struct {
	nodeID      string
	description string
}

type DotGraphOptions map[string]string

//New creates a new Graph with a given name
func New(name string) *DotGraph {
	return &DotGraph{
		name:         name,
		edges:        make(map[string][]edge),
		edgeOptions:  make(map[string]DotGraphOptions),
		nodeOptions:  DotGraphOptions{},
		graphOptions: DotGraphOptions{},
	}
}

//AddNode add a node to string. There does not have to be an edge for a node.
func (g DotGraph) AddNode(node string) {
	new := getIDSafeNodeName(node)
	if g.edges[new] == nil {
		g.edges[new] = []edge{}
	}
}

//GetDotFileContent create the content of a dot-file (graphviz)
func (g DotGraph) String() string {
	content := []string{"digraph " + g.name + " {"}
	if graphStyle := createGraphOptionsString(g.graphOptions); graphStyle != "" {
		content = append(content, "graph "+graphStyle)
	}
	if nodeStyle := createGraphOptionsString(g.nodeOptions); nodeStyle != "" {
		content = append(content, "node "+nodeStyle)
	}

	for from, deps := range g.edges {
		content = append(content, from)

		for _, to := range deps {
			if from != `""` && to.nodeID != `""` {
				edgeStyle := g.createEdgeOptionsString(from, to.nodeID, to.description)
				content = append(content, from+"->"+to.nodeID+edgeStyle)
			}
		}
	}
	content = append(content, "}")

	return strings.Join(content, "\n")
}

func (g DotGraph) createEdgeOptionsString(from, to, description string) string {
	options := []string{}
	if description != "" {
		options = append(options, "label=\""+description+"\"")
	}

	for pattern, ops := range g.edgeOptions {
		matchesFrom, _ := regexp.MatchString(pattern, from)
		matchesTo, _ := regexp.MatchString(pattern, to)

		if matchesFrom || matchesTo {
			for attr, val := range ops {
				options = append(options, strings.ToLower(attr)+"=\""+val+"\"")
			}
		}
	}
	if len(options) == 0 {
		return ""
	}

	return "[" + strings.Join(options, " ") + "]"
}

func createGraphOptionsString(gOptions DotGraphOptions) string {
	options := []string{}
	for attr, val := range gOptions {
		options = append(options, strings.ToLower(attr)+"=\""+val+"\"")
	}
	if len(options) == 0 {
		return ""
	}

	return "[" + strings.Join(options, " ") + "]"
}

//AddEdgeGraphOptions adds options that applies options to edges when pattern matches
func (g DotGraph) AddEdgeGraphOptions(pattern string, options DotGraphOptions) {
	g.edgeOptions[pattern] = options
}

//SetNodeGraphOptions sets options that apply to all nodes
func (g DotGraph) SetNodeGraphOptions(options DotGraphOptions) {
	for k, v := range options {
		g.nodeOptions[k] = v
	}
}

//SetGraphOptions sets options that apply to the graph
func (g DotGraph) SetGraphOptions(options DotGraphOptions) {
	for k, v := range options {
		g.graphOptions[k] = v
	}
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
func (g DotGraph) AddDirectedEdge(from, to, description string) {
	saveFrom := getIDSafeNodeName(from)
	saveTo := getIDSafeNodeName(to)

	if _, found := g.edges[saveFrom]; !found {
		g.edges[saveFrom] = []edge{}
	}

	for _, edge := range g.edges[saveFrom] {
		if edge.nodeID == saveTo {
			return
		}
	}

	g.edges[saveFrom] = append(g.edges[saveFrom], edge{
		nodeID:      saveTo,
		description: description,
	})
}

//GetDependencies returns alls direct dipendencies for a package within the graph
func (g DotGraph) GetDependencies(pkg string) []string {
	dependencies := []string{}

	for from, deps := range g.edges {
		if from == getIDSafeNodeName(pkg) {
			for _, edge := range deps {
				dependencies = append(dependencies, edge.nodeID)
			}
		}
	}
	return dependencies
}

//GetDependents returns all packages that directly depend on the given package within the graph
func (g DotGraph) GetDependents(pkg string) []string {
	dependents := []string{}
loop:
	for from, deps := range g.edges {
		for _, to := range deps {
			if to.nodeID == getIDSafeNodeName(pkg) {
				dependents = append(dependents, from)
				continue loop
			}
		}
	}
	return dependents
}
