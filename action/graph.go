package action

type Graph interface {
	AddNode(node string)
	AddDirectedEdge(from, to string)
	GetDependencies(pkg string) []string
	GetDependents(pkg string) []string
	String() string
}

type GraphRenderer interface {
	Render(graphContent string)
}
