package graph

import "sort"

// StronglyConnectedComponents returns the maximal strongly connected
// components of the snapshot.
func (s *Snapshot[NodeID, EdgeID, NodeData, EdgeData]) StronglyConnectedComponents() [][]NodeID {
	return StronglyConnectedComponents(s)
}

// StronglyConnectedComponents returns the maximal strongly connected components of g.
//
// Nodes within each component and the components themselves are ordered by the
// stable order supplied by NodeIDs.
func StronglyConnectedComponents[NodeID comparable, EdgeID comparable, EdgeData any](
	g FiniteDirected[NodeID, EdgeID, EdgeData],
) [][]NodeID {
	nodeIDs := g.NodeIDs()
	if nodeIDs == nil {
		return [][]NodeID{}
	}

	nodes := make([]NodeID, 0)
	position := make(map[NodeID]int)
	for node := range nodeIDs {
		position[node] = len(nodes)
		nodes = append(nodes, node)
	}

	nextIndex := 0
	indices := make(map[NodeID]int)
	lowlink := make(map[NodeID]int)
	onStack := make(map[NodeID]bool)
	stack := make([]NodeID, 0)
	components := make([][]NodeID, 0)

	var visit func(NodeID)
	visit = func(node NodeID) {
		indices[node] = nextIndex
		lowlink[node] = nextIndex
		nextIndex++
		stack = append(stack, node)
		onStack[node] = true

		edges := g.Outgoing(node)
		if edges != nil {
			for edge := range edges {
				index, seen := indices[edge.To]
				if !seen {
					visit(edge.To)
					lowlink[node] = min(lowlink[node], lowlink[edge.To])
					continue
				}
				if onStack[edge.To] {
					lowlink[node] = min(lowlink[node], index)
				}
			}
		}

		if lowlink[node] != indices[node] {
			return
		}

		component := make([]NodeID, 0)
		for {
			last := len(stack) - 1
			member := stack[last]
			stack = stack[:last]
			onStack[member] = false
			component = append(component, member)
			if member == node {
				break
			}
		}
		components = append(components, component)
	}

	for _, node := range nodes {
		if _, seen := indices[node]; !seen {
			visit(node)
		}
	}

	for _, component := range components {
		sort.Slice(component, func(i, j int) bool {
			return position[component[i]] < position[component[j]]
		})
	}
	sort.Slice(components, func(i, j int) bool {
		return position[components[i][0]] < position[components[j][0]]
	})

	return components
}
