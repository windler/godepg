package action

// GraphFilter provides filters that are applied before or after a node will be added to the graph
type GraphFilter interface {
	GetPreNodeFilters() []Matcher
	GetPostNodeFilters() []Matcher
}

//Context provides flags
type Context interface {
	GetStringFlag(flag string) string
	GetStringSliceFlag(flag string) []string
	GetIntFlag(flag string) int
	GetBoolFlag(flag string) bool
}

//Matcher filters nodes
type Matcher interface {
	Matches() bool
}
