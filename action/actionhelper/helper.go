package actionhelper

import "github.com/windler/godepg/action"

//AddEdge adds a edge to the graph applying filters
func AddEdge(graph action.Graph, from, to, description string, filter action.GraphFilter) bool {
	for _, f := range filter.GetPreNodeFilters() {
		if f.Matches() {
			return false
		}
	}

	graph.AddNode(from)

	for _, f := range filter.GetPostNodeFilters() {
		if f.Matches() {
			return false
		}
	}

	graph.AddDirectedEdge(from, to, description)
	return true
}
