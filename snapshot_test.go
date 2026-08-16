package graph

import (
	"errors"
	"iter"
	"reflect"
	"slices"
	"testing"
)

type testSource struct {
	nodes []Node[string, string]
	edges []Edge[string, string, string]
}

type nilSource struct{}

func (nilSource) Nodes() iter.Seq[Node[string, string]] {
	return nil
}

func (nilSource) Edges() iter.Seq[Edge[string, string, string]] {
	return nil
}

func (s testSource) Nodes() iter.Seq[Node[string, string]] {
	return slices.Values(s.nodes)
}

func (s testSource) Edges() iter.Seq[Edge[string, string, string]] {
	return slices.Values(s.edges)
}

func TestBuildSnapshot(t *testing.T) {
	source := testSource{
		nodes: []Node[string, string]{
			{ID: "a", Data: "node-a"},
			{ID: "b", Data: "node-b"},
			{ID: "c", Data: "node-c"},
		},
		edges: []Edge[string, string, string]{
			{ID: "supports", From: "a", To: "b", Data: "first"},
			{ID: "depends", From: "a", To: "b", Data: "second"},
			{ID: "self", From: "c", To: "c", Data: "third"},
		},
	}

	snapshot, err := Build(source)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	if !snapshot.HasNode("b") {
		t.Fatal("HasNode(\"b\") = false, want true")
	}
	if snapshot.HasNode("missing") {
		t.Fatal("HasNode(\"missing\") = true, want false")
	}

	node, ok := snapshot.Node("b")
	if !ok || node.Data != "node-b" {
		t.Fatalf("Node(\"b\") = %#v, %t", node, ok)
	}
	edge, ok := snapshot.Edge("depends")
	if !ok || edge.Data != "second" {
		t.Fatalf("Edge(\"depends\") = %#v, %t", edge, ok)
	}

	if got, want := slices.Collect(snapshot.NodeIDs()), []string{"a", "b", "c"}; !reflect.DeepEqual(got, want) {
		t.Errorf("NodeIDs() = %v, want %v", got, want)
	}
	if got, want := edgeIDs(snapshot.Edges()), []string{"supports", "depends", "self"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Edges() IDs = %v, want %v", got, want)
	}
	if got, want := edgeIDs(snapshot.Outgoing("a")), []string{"supports", "depends"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Outgoing(\"a\") IDs = %v, want %v", got, want)
	}
	if got, want := edgeIDs(snapshot.Incoming("b")), []string{"supports", "depends"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Incoming(\"b\") IDs = %v, want %v", got, want)
	}
}

func TestBuildCopiesTopology(t *testing.T) {
	source := testSource{
		nodes: []Node[string, string]{{ID: "a"}, {ID: "b"}},
		edges: []Edge[string, string, string]{{ID: "ab", From: "a", To: "b"}},
	}

	snapshot, err := Build(source)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	source.nodes[0].ID = "changed"
	source.edges[0].To = "changed"

	if got, want := slices.Collect(snapshot.NodeIDs()), []string{"a", "b"}; !reflect.DeepEqual(got, want) {
		t.Errorf("NodeIDs() after source mutation = %v, want %v", got, want)
	}
	edge, ok := snapshot.Edge("ab")
	if !ok || edge.To != "b" {
		t.Errorf("Edge(\"ab\") after source mutation = %#v, %t", edge, ok)
	}
}

func TestBuildEmptySnapshot(t *testing.T) {
	tests := []struct {
		name   string
		source Source[string, string, string, string]
	}{
		{name: "nil iterators", source: nilSource{}},
		{name: "empty iterators", source: testSource{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snapshot, err := Build(tt.source)
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}
			if got := slices.Collect(snapshot.NodeIDs()); len(got) != 0 {
				t.Errorf("NodeIDs() = %v, want empty", got)
			}
			if got := slices.Collect(snapshot.Edges()); len(got) != 0 {
				t.Errorf("Edges() = %v, want empty", got)
			}
		})
	}
}

func TestBuildRejectsInvalidTopology(t *testing.T) {
	tests := []struct {
		name    string
		source  Source[string, string, string, string]
		wantErr error
	}{
		{
			name: "edge without nodes",
			source: testSource{
				edges: []Edge[string, string, string]{{ID: "edge", From: "a", To: "b"}},
			},
			wantErr: ErrMissingEndpoint,
		},
		{
			name: "duplicate node",
			source: testSource{
				nodes: []Node[string, string]{{ID: "a"}, {ID: "a"}},
			},
			wantErr: ErrDuplicateNode,
		},
		{
			name: "duplicate edge",
			source: testSource{
				nodes: []Node[string, string]{{ID: "a"}},
				edges: []Edge[string, string, string]{
					{ID: "edge", From: "a", To: "a"},
					{ID: "edge", From: "a", To: "a"},
				},
			},
			wantErr: ErrDuplicateEdge,
		},
		{
			name: "missing from endpoint",
			source: testSource{
				nodes: []Node[string, string]{{ID: "a"}},
				edges: []Edge[string, string, string]{{ID: "edge", From: "missing", To: "a"}},
			},
			wantErr: ErrMissingEndpoint,
		},
		{
			name: "missing to endpoint",
			source: testSource{
				nodes: []Node[string, string]{{ID: "a"}},
				edges: []Edge[string, string, string]{{ID: "edge", From: "a", To: "missing"}},
			},
			wantErr: ErrMissingEndpoint,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Build(tt.source)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Build() error = %v, want errors.Is(%v)", err, tt.wantErr)
			}
		})
	}
}

func edgeIDs(edges iter.Seq[Edge[string, string, string]]) []string {
	ids := make([]string, 0)
	for edge := range edges {
		ids = append(ids, edge.ID)
	}
	return ids
}
