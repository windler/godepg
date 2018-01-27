package action

//Graph represents the dependy graph
type Graph interface {
	AddNode(node string)
	AddDirectedEdge(from, to, description string)
	GetDependencies(pkg string) []string
	GetDependents(pkg string) []string
	String() string
}

//GraphRenderer renders the graph
type GraphRenderer interface {
	Render(graphContent string)
}
