package graph

// CycleError reports that a topological ordering does not exist and includes
// one structural cycle witness.
type CycleError[NodeID comparable, EdgeID comparable, EdgeData any] struct {
	Cycle Cycle[NodeID, EdgeID, EdgeData]
}

// Error implements error.
func (e *CycleError[NodeID, EdgeID, EdgeData]) Error() string {
	return "graph contains a cycle"
}

// TopologicalSort returns a deterministic topological ordering of g.
//
// Zero-indegree nodes are initially queued in NodeIDs order. Nodes made ready
// during traversal are queued in outgoing edge order. On a cyclic graph, the
// function returns a CycleError with a witness and no partial ordering. It
// returns ErrInconsistentGraph if Kahn's algorithm and cycle detection disagree.
func TopologicalSort[NodeID comparable, EdgeID comparable, EdgeData any](
	g FiniteDirected[NodeID, EdgeID, EdgeData],
) ([]NodeID, error) {
	nodeIDs := g.NodeIDs()
	if nodeIDs == nil {
		return []NodeID{}, nil
	}

	nodes := make([]NodeID, 0)
	indegree := make(map[NodeID]int)
	for node := range nodeIDs {
		nodes = append(nodes, node)
		indegree[node] = 0
	}

	for _, node := range nodes {
		edges := g.Outgoing(node)
		if edges == nil {
			continue
		}
		for edge := range edges {
			indegree[edge.To]++
		}
	}

	queue := make([]NodeID, 0)
	for _, node := range nodes {
		if indegree[node] == 0 {
			queue = append(queue, node)
		}
	}

	order := make([]NodeID, 0, len(nodes))
	for head := 0; head < len(queue); head++ {
		node := queue[head]
		order = append(order, node)

		edges := g.Outgoing(node)
		if edges == nil {
			continue
		}
		for edge := range edges {
			indegree[edge.To]--
			if indegree[edge.To] == 0 {
				queue = append(queue, edge.To)
			}
		}
	}

	if len(order) == len(nodes) {
		return order, nil
	}

	cycle, found := FindCycle(g)
	if !found {
		return nil, ErrInconsistentGraph
	}
	return nil, &CycleError[NodeID, EdgeID, EdgeData]{Cycle: cycle}
}
