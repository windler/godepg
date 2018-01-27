package actionhelper

import (
	"testing"

	"github.com/windler/godepg/action"

	"github.com/windler/godepg/action/mocks"
)

func TestAddEdge(t *testing.T) {
	graphMock := new(mocks.Graph)
	filterMock := new(mocks.GraphFilter)

	filterMock.On("GetPreNodeFilters").Return([]action.Matcher{}).Once()
	filterMock.On("GetPostNodeFilters").Return([]action.Matcher{}).Once()

	graphMock.On("AddDirectedEdge", "A", "B", "description").Once()
	graphMock.On("AddNode", "A").Once()

	AddEdge(graphMock, "A", "B", "description", filterMock)

	filterMock.AssertExpectations(t)
	graphMock.AssertExpectations(t)
}

func TestAddEdgePreFilter(t *testing.T) {
	graphMock := new(mocks.Graph)
	filterMock := new(mocks.GraphFilter)
	matcherMock := new(mocks.Matcher)

	filterMock.On("GetPreNodeFilters").Return([]action.Matcher{matcherMock}).Once()

	matcherMock.On("Matches").Return(true).Once()

	AddEdge(graphMock, "A", "B", "description", filterMock)

	filterMock.AssertExpectations(t)
	graphMock.AssertExpectations(t)
	matcherMock.AssertExpectations(t)
}

func TestAddEdgePostFilte(t *testing.T) {
	graphMock := new(mocks.Graph)
	filterMock := new(mocks.GraphFilter)
	matcherMock := new(mocks.Matcher)

	filterMock.On("GetPreNodeFilters").Return([]action.Matcher{}).Once()
	filterMock.On("GetPostNodeFilters").Return([]action.Matcher{matcherMock}).Once()

	matcherMock.On("Matches").Return(true).Once()

	graphMock.On("AddNode", "A").Once()

	AddEdge(graphMock, "A", "B", "description", filterMock)

	filterMock.AssertExpectations(t)
	graphMock.AssertExpectations(t)
	matcherMock.AssertExpectations(t)
}
