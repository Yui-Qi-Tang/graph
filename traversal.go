package graph

import "iter"

// BFS traverses nodes breadth-first from start.
//
// The traversal yields start first and then follows outgoing edges in the
// order supplied by g. Each node is yielded at most once. Callers must provide
// a start node that belongs to the graph.
func BFS[NodeID comparable, EdgeID comparable, EdgeData any](
	g Directed[NodeID, EdgeID, EdgeData],
	start NodeID,
) iter.Seq[NodeID] {
	return func(yield func(NodeID) bool) {
		seen := map[NodeID]struct{}{start: {}}
		queue := []NodeID{start}

		for head := 0; head < len(queue); head++ {
			node := queue[head]
			if !yield(node) {
				return
			}

			edges := g.Outgoing(node)
			if edges == nil {
				continue
			}
			for edge := range edges {
				if _, exists := seen[edge.To]; exists {
					continue
				}
				seen[edge.To] = struct{}{}
				queue = append(queue, edge.To)
			}
		}
	}
}

// DFS traverses nodes depth-first in pre-order from start.
//
// Outgoing edge order determines which branch is visited first. Each node is
// yielded at most once. Callers must provide a start node that belongs to the graph.
func DFS[NodeID comparable, EdgeID comparable, EdgeData any](
	g Directed[NodeID, EdgeID, EdgeData],
	start NodeID,
) iter.Seq[NodeID] {
	return func(yield func(NodeID) bool) {
		seen := make(map[NodeID]struct{})
		stack := []NodeID{start}

		for len(stack) > 0 {
			last := len(stack) - 1
			node := stack[last]
			stack = stack[:last]

			if _, exists := seen[node]; exists {
				continue
			}
			seen[node] = struct{}{}
			if !yield(node) {
				return
			}

			edges := g.Outgoing(node)
			if edges == nil {
				continue
			}
			next := make([]NodeID, 0)
			for edge := range edges {
				next = append(next, edge.To)
			}
			for i := len(next) - 1; i >= 0; i-- {
				stack = append(stack, next[i])
			}
		}
	}
}

// Reachable reports whether to can be reached from from.
// A node is reachable from itself without traversing an edge.
func Reachable[NodeID comparable, EdgeID comparable, EdgeData any](
	g Directed[NodeID, EdgeID, EdgeData],
	from NodeID,
	to NodeID,
) bool {
	for node := range BFS(g, from) {
		if node == to {
			return true
		}
	}
	return false
}
