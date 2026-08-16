package graph

import "slices"

// Path is a directed path whose edge records retain caller-defined data.
// Nodes contains both endpoints, so a non-empty edge path has one more node
// than edge. A zero-edge path from a node to itself contains that node once.
type Path[NodeID comparable, EdgeID comparable, EdgeData any] struct {
	Nodes []NodeID
	Edges []Edge[NodeID, EdgeID, EdgeData]
}

// FindPath finds the first shortest path by edge count from from to to.
// Outgoing edge order breaks ties deterministically.
func FindPath[NodeID comparable, EdgeID comparable, EdgeData any](
	g Directed[NodeID, EdgeID, EdgeData],
	from NodeID,
	to NodeID,
) (Path[NodeID, EdgeID, EdgeData], bool) {
	if from == to {
		return Path[NodeID, EdgeID, EdgeData]{Nodes: []NodeID{from}}, true
	}

	seen := map[NodeID]struct{}{from: {}}
	previous := make(map[NodeID]Edge[NodeID, EdgeID, EdgeData])
	queue := []NodeID{from}

	for head := 0; head < len(queue); head++ {
		node := queue[head]
		edges := g.Outgoing(node)
		if edges == nil {
			continue
		}
		for edge := range edges {
			if _, exists := seen[edge.To]; exists {
				continue
			}
			seen[edge.To] = struct{}{}
			previous[edge.To] = edge
			if edge.To == to {
				return buildPath(previous, from, to), true
			}
			queue = append(queue, edge.To)
		}
	}

	return Path[NodeID, EdgeID, EdgeData]{}, false
}

func buildPath[NodeID comparable, EdgeID comparable, EdgeData any](
	previous map[NodeID]Edge[NodeID, EdgeID, EdgeData],
	from NodeID,
	to NodeID,
) Path[NodeID, EdgeID, EdgeData] {
	nodes := []NodeID{to}
	edges := make([]Edge[NodeID, EdgeID, EdgeData], 0)

	for node := to; node != from; {
		edge := previous[node]
		edges = append(edges, edge)
		node = edge.From
		nodes = append(nodes, node)
	}

	slices.Reverse(nodes)
	slices.Reverse(edges)
	return Path[NodeID, EdgeID, EdgeData]{Nodes: nodes, Edges: edges}
}
