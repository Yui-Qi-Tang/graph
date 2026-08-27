package graph

import "slices"

// Cycle is a directed cycle witness with caller-defined edge data.
// Nodes repeats the starting node at the end, so len(Nodes) is len(Edges)+1.
type Cycle[NodeID comparable, EdgeID comparable, EdgeData any] struct {
	Nodes []NodeID
	Edges []Edge[NodeID, EdgeID, EdgeData]
}

// HasCycle reports whether the snapshot contains a directed cycle.
func (s *Snapshot[NodeID, EdgeID, NodeData, EdgeData]) HasCycle() bool {
	return HasCycle(s)
}

// HasCycle reports whether g contains a directed cycle.
func HasCycle[NodeID comparable, EdgeID comparable, EdgeData any](
	g FiniteDirected[NodeID, EdgeID, EdgeData],
) bool {
	_, found := FindCycle(g)
	return found
}

// FindCycle returns the first directed cycle found in stable graph order.
func (s *Snapshot[NodeID, EdgeID, NodeData, EdgeData]) FindCycle() (Cycle[NodeID, EdgeID, EdgeData], bool) {
	return FindCycle(s)
}

// FindCycle returns the first directed cycle found in stable graph order.
// The result is a structural witness, not a domain-level interpretation.
func FindCycle[NodeID comparable, EdgeID comparable, EdgeData any](
	g FiniteDirected[NodeID, EdgeID, EdgeData],
) (Cycle[NodeID, EdgeID, EdgeData], bool) {
	const (
		unvisited uint8 = iota
		visiting
		visited
	)

	state := make(map[NodeID]uint8)
	parent := make(map[NodeID]Edge[NodeID, EdgeID, EdgeData])
	var found Cycle[NodeID, EdgeID, EdgeData]

	var visit func(NodeID) bool
	visit = func(node NodeID) bool {
		state[node] = visiting
		edges := g.Outgoing(node)
		if edges != nil {
			for edge := range edges {
				switch state[edge.To] {
				case unvisited:
					parent[edge.To] = edge
					if visit(edge.To) {
						return true
					}
				case visiting:
					found = cycleFromBackEdge(parent, node, edge)
					return true
				}
			}
		}
		state[node] = visited
		return false
	}

	nodes := g.NodeIDs()
	if nodes == nil {
		return Cycle[NodeID, EdgeID, EdgeData]{}, false
	}
	for node := range nodes {
		if state[node] != unvisited {
			continue
		}
		if visit(node) {
			return found, true
		}
	}

	return Cycle[NodeID, EdgeID, EdgeData]{}, false
}

func cycleFromBackEdge[NodeID comparable, EdgeID comparable, EdgeData any](
	parent map[NodeID]Edge[NodeID, EdgeID, EdgeData],
	from NodeID,
	back Edge[NodeID, EdgeID, EdgeData],
) Cycle[NodeID, EdgeID, EdgeData] {
	nodes := []NodeID{from}
	edges := make([]Edge[NodeID, EdgeID, EdgeData], 0)

	for node := from; node != back.To; {
		edge := parent[node]
		edges = append(edges, edge)
		node = edge.From
		nodes = append(nodes, node)
	}

	slices.Reverse(nodes)
	slices.Reverse(edges)
	nodes = append(nodes, back.To)
	edges = append(edges, back)

	return Cycle[NodeID, EdgeID, EdgeData]{Nodes: nodes, Edges: edges}
}
