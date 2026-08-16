package graph

import (
	"fmt"
	"iter"
)

// Snapshot is an immutable directed graph topology built from a Source.
//
// Snapshot copies node and edge records, preserves their iteration order, and
// builds private incoming and outgoing indexes. The topology cannot be mutated
// after construction. Node and edge data are shallowly copied.
type Snapshot[NodeID comparable, EdgeID comparable, NodeData, EdgeData any] struct {
	nodes map[NodeID]Node[NodeID, NodeData]
	edges map[EdgeID]Edge[NodeID, EdgeID, EdgeData]

	nodeOrder []NodeID
	edgeOrder []EdgeID
	outgoing  map[NodeID][]EdgeID
	incoming  map[NodeID][]EdgeID
}

// Build validates source topology and returns an immutable Snapshot.
//
// An empty source is valid. Node and edge IDs must be unique, and every edge
// endpoint must reference a yielded node. Self-loops and parallel edges are
// valid.
func Build[NodeID comparable, EdgeID comparable, NodeData, EdgeData any](
	source Source[NodeID, EdgeID, NodeData, EdgeData],
) (*Snapshot[NodeID, EdgeID, NodeData, EdgeData], error) {
	s := &Snapshot[NodeID, EdgeID, NodeData, EdgeData]{
		nodes:    make(map[NodeID]Node[NodeID, NodeData]),
		edges:    make(map[EdgeID]Edge[NodeID, EdgeID, EdgeData]),
		outgoing: make(map[NodeID][]EdgeID),
		incoming: make(map[NodeID][]EdgeID),
	}

	nodes := source.Nodes()
	if nodes != nil {
		for node := range nodes {
			if _, exists := s.nodes[node.ID]; exists {
				return nil, fmt.Errorf("%w: %v", ErrDuplicateNode, node.ID)
			}
			s.nodes[node.ID] = node
			s.nodeOrder = append(s.nodeOrder, node.ID)
		}
	}

	edges := source.Edges()
	if edges != nil {
		for edge := range edges {
			if _, exists := s.edges[edge.ID]; exists {
				return nil, fmt.Errorf("%w: %v", ErrDuplicateEdge, edge.ID)
			}
			if _, exists := s.nodes[edge.From]; !exists {
				return nil, fmt.Errorf("%w: edge %v from %v", ErrMissingEndpoint, edge.ID, edge.From)
			}
			if _, exists := s.nodes[edge.To]; !exists {
				return nil, fmt.Errorf("%w: edge %v to %v", ErrMissingEndpoint, edge.ID, edge.To)
			}

			s.edges[edge.ID] = edge
			s.edgeOrder = append(s.edgeOrder, edge.ID)
			s.outgoing[edge.From] = append(s.outgoing[edge.From], edge.ID)
			s.incoming[edge.To] = append(s.incoming[edge.To], edge.ID)
		}
	}

	return s, nil
}

// HasNode reports whether id belongs to the snapshot.
func (s *Snapshot[NodeID, EdgeID, NodeData, EdgeData]) HasNode(id NodeID) bool {
	_, ok := s.nodes[id]
	return ok
}

// Node returns the node identified by id.
func (s *Snapshot[NodeID, EdgeID, NodeData, EdgeData]) Node(id NodeID) (Node[NodeID, NodeData], bool) {
	node, ok := s.nodes[id]
	return node, ok
}

// Edge returns the edge identified by id.
func (s *Snapshot[NodeID, EdgeID, NodeData, EdgeData]) Edge(id EdgeID) (Edge[NodeID, EdgeID, EdgeData], bool) {
	edge, ok := s.edges[id]
	return edge, ok
}

// NodeIDs yields node IDs in source order.
func (s *Snapshot[NodeID, EdgeID, NodeData, EdgeData]) NodeIDs() iter.Seq[NodeID] {
	return func(yield func(NodeID) bool) {
		for _, id := range s.nodeOrder {
			if !yield(id) {
				return
			}
		}
	}
}

// Nodes yields nodes in source order.
func (s *Snapshot[NodeID, EdgeID, NodeData, EdgeData]) Nodes() iter.Seq[Node[NodeID, NodeData]] {
	return func(yield func(Node[NodeID, NodeData]) bool) {
		for _, id := range s.nodeOrder {
			if !yield(s.nodes[id]) {
				return
			}
		}
	}
}

// Edges yields edges in source order.
func (s *Snapshot[NodeID, EdgeID, NodeData, EdgeData]) Edges() iter.Seq[Edge[NodeID, EdgeID, EdgeData]] {
	return func(yield func(Edge[NodeID, EdgeID, EdgeData]) bool) {
		for _, id := range s.edgeOrder {
			if !yield(s.edges[id]) {
				return
			}
		}
	}
}

// Outgoing yields edges leaving id in source edge order.
func (s *Snapshot[NodeID, EdgeID, NodeData, EdgeData]) Outgoing(id NodeID) iter.Seq[Edge[NodeID, EdgeID, EdgeData]] {
	return func(yield func(Edge[NodeID, EdgeID, EdgeData]) bool) {
		for _, edgeID := range s.outgoing[id] {
			if !yield(s.edges[edgeID]) {
				return
			}
		}
	}
}

// Incoming yields edges entering id in source edge order.
func (s *Snapshot[NodeID, EdgeID, NodeData, EdgeData]) Incoming(id NodeID) iter.Seq[Edge[NodeID, EdgeID, EdgeData]] {
	return func(yield func(Edge[NodeID, EdgeID, EdgeData]) bool) {
		for _, edgeID := range s.incoming[id] {
			if !yield(s.edges[edgeID]) {
				return
			}
		}
	}
}
