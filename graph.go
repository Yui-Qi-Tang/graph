package graph

import (
	"errors"
	"iter"
)

var (
	// ErrDuplicateNode reports that a source yielded the same node ID twice.
	ErrDuplicateNode = errors.New("duplicate node")
	// ErrDuplicateEdge reports that a source yielded the same edge ID twice.
	ErrDuplicateEdge = errors.New("duplicate edge")
	// ErrMissingEndpoint reports that an edge references a node outside the source.
	ErrMissingEndpoint = errors.New("missing edge endpoint")
	// ErrInconsistentGraph reports that a graph violated the FiniteDirected contract.
	ErrInconsistentGraph = errors.New("inconsistent graph")
)

// Node is a graph node with caller-defined data.
//
// Data is copied by assignment when a Snapshot is built. Values containing
// pointers, slices, maps, or other reference-like fields therefore remain
// shallowly shared with the caller.
type Node[ID comparable, Data any] struct {
	ID   ID
	Data Data
}

// Edge is a directed graph edge with an independent identity and caller-defined data.
//
// Independent edge IDs allow parallel edges between the same pair of nodes.
// Data has the same shallow-copy semantics as Node.Data.
type Edge[NodeID comparable, EdgeID comparable, Data any] struct {
	ID   EdgeID
	From NodeID
	To   NodeID
	Data Data
}

// Source supplies the complete node and edge records used to build a Snapshot.
// Iteration order is preserved by the snapshot and its adjacency indexes.
type Source[NodeID comparable, EdgeID comparable, NodeData, EdgeData any] interface {
	Nodes() iter.Seq[Node[NodeID, NodeData]]
	Edges() iter.Seq[Edge[NodeID, EdgeID, EdgeData]]
}

// Directed supplies outgoing edges for structural algorithms.
//
// Outgoing must be stable during an algorithm call. Every yielded edge must
// have From equal to the requested node.
type Directed[NodeID comparable, EdgeID comparable, EdgeData any] interface {
	Outgoing(NodeID) iter.Seq[Edge[NodeID, EdgeID, EdgeData]]
}

// FiniteDirected is a complete finite directed graph.
//
// NodeIDs must yield each node exactly once in stable order. Every endpoint
// yielded by Outgoing must also be present in NodeIDs.
type FiniteDirected[NodeID comparable, EdgeID comparable, EdgeData any] interface {
	Directed[NodeID, EdgeID, EdgeData]
	NodeIDs() iter.Seq[NodeID]
}
