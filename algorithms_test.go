package graph

import (
	"errors"
	"iter"
	"reflect"
	"slices"
	"testing"
)

type inconsistentGraph struct{}

func (inconsistentGraph) NodeIDs() iter.Seq[string] {
	return slices.Values([]string{"a"})
}

func (inconsistentGraph) Outgoing(node string) iter.Seq[Edge[string, string, string]] {
	if node != "a" {
		return nil
	}
	return slices.Values([]Edge[string, string, string]{
		{ID: "missing-endpoint", From: "a", To: "missing"},
	})
}

func TestTraversal(t *testing.T) {
	g := buildTestSnapshot(t,
		[]string{"a", "b", "c", "d", "e"},
		[]Edge[string, string, string]{
			{ID: "ab", From: "a", To: "b"},
			{ID: "ac", From: "a", To: "c"},
			{ID: "bd", From: "b", To: "d"},
			{ID: "cd", From: "c", To: "d"},
			{ID: "de", From: "d", To: "e"},
		},
	)

	if got, want := slices.Collect(BFS(g, "a")), []string{"a", "b", "c", "d", "e"}; !reflect.DeepEqual(got, want) {
		t.Errorf("BFS() = %v, want %v", got, want)
	}
	if got, want := slices.Collect(DFS(g, "a")), []string{"a", "b", "d", "e", "c"}; !reflect.DeepEqual(got, want) {
		t.Errorf("DFS() = %v, want %v", got, want)
	}
	if !Reachable(g, "a", "e") {
		t.Error("Reachable(a, e) = false, want true")
	}
	if Reachable(g, "e", "a") {
		t.Error("Reachable(e, a) = true, want false")
	}
	if !Reachable(g, "c", "c") {
		t.Error("Reachable(c, c) = false, want true")
	}
}

func TestFindPath(t *testing.T) {
	g := buildTestSnapshot(t,
		[]string{"a", "b", "c", "d", "e"},
		[]Edge[string, string, string]{
			{ID: "ab", From: "a", To: "b", Data: "first"},
			{ID: "ac", From: "a", To: "c", Data: "second"},
			{ID: "bd", From: "b", To: "d", Data: "third"},
			{ID: "cd", From: "c", To: "d", Data: "fourth"},
			{ID: "de", From: "d", To: "e", Data: "fifth"},
		},
	)

	path, found := FindPath(g, "a", "e")
	if !found {
		t.Fatal("FindPath(a, e) found = false, want true")
	}
	if got, want := path.Nodes, []string{"a", "b", "d", "e"}; !reflect.DeepEqual(got, want) {
		t.Errorf("FindPath(a, e) nodes = %v, want %v", got, want)
	}
	if got, want := edgeIDs(slices.Values(path.Edges)), []string{"ab", "bd", "de"}; !reflect.DeepEqual(got, want) {
		t.Errorf("FindPath(a, e) edge IDs = %v, want %v", got, want)
	}
	if path.Edges[0].Data != "first" {
		t.Errorf("FindPath(a, e) first edge data = %q, want %q", path.Edges[0].Data, "first")
	}

	same, found := FindPath(g, "c", "c")
	if !found || !reflect.DeepEqual(same.Nodes, []string{"c"}) || len(same.Edges) != 0 {
		t.Errorf("FindPath(c, c) = %#v, %t", same, found)
	}
	if _, found := FindPath(g, "e", "a"); found {
		t.Error("FindPath(e, a) found = true, want false")
	}
}

func TestFindCycle(t *testing.T) {
	g := buildTestSnapshot(t,
		[]string{"a", "b", "c", "d"},
		[]Edge[string, string, string]{
			{ID: "ab", From: "a", To: "b", Data: "first"},
			{ID: "bc", From: "b", To: "c", Data: "second"},
			{ID: "ca", From: "c", To: "a", Data: "third"},
			{ID: "dd", From: "d", To: "d", Data: "self"},
		},
	)

	cycle, found := FindCycle(g)
	if !found {
		t.Fatal("FindCycle() found = false, want true")
	}
	if got, want := cycle.Nodes, []string{"a", "b", "c", "a"}; !reflect.DeepEqual(got, want) {
		t.Errorf("FindCycle() nodes = %v, want %v", got, want)
	}
	if got, want := edgeIDs(slices.Values(cycle.Edges)), []string{"ab", "bc", "ca"}; !reflect.DeepEqual(got, want) {
		t.Errorf("FindCycle() edge IDs = %v, want %v", got, want)
	}
	if cycle.Edges[2].Data != "third" {
		t.Errorf("FindCycle() last edge data = %q, want %q", cycle.Edges[2].Data, "third")
	}
	if !HasCycle(g) {
		t.Error("HasCycle() = false, want true")
	}
}

func TestFindCycleSelfLoop(t *testing.T) {
	g := buildTestSnapshot(t,
		[]string{"a"},
		[]Edge[string, string, string]{{ID: "aa", From: "a", To: "a"}},
	)

	cycle, found := FindCycle(g)
	if !found {
		t.Fatal("FindCycle() found = false, want true")
	}
	if got, want := cycle.Nodes, []string{"a", "a"}; !reflect.DeepEqual(got, want) {
		t.Errorf("FindCycle() nodes = %v, want %v", got, want)
	}
}

func TestTopologicalSort(t *testing.T) {
	g := buildTestSnapshot(t,
		[]string{"a", "b", "c", "d"},
		[]Edge[string, string, string]{
			{ID: "ac-1", From: "a", To: "c"},
			{ID: "ac-2", From: "a", To: "c"},
			{ID: "bc", From: "b", To: "c"},
			{ID: "cd", From: "c", To: "d"},
		},
	)

	order, err := TopologicalSort(g)
	if err != nil {
		t.Fatalf("TopologicalSort() error = %v", err)
	}
	if want := []string{"a", "b", "c", "d"}; !reflect.DeepEqual(order, want) {
		t.Errorf("TopologicalSort() = %v, want %v", order, want)
	}
}

func TestTopologicalSortCycleError(t *testing.T) {
	g := buildTestSnapshot(t,
		[]string{"a", "b"},
		[]Edge[string, string, string]{
			{ID: "ab", From: "a", To: "b"},
			{ID: "ba", From: "b", To: "a"},
		},
	)

	order, err := TopologicalSort(g)
	if err == nil {
		t.Fatal("TopologicalSort() error = nil, want cycle error")
	}
	if order != nil {
		t.Errorf("TopologicalSort() order = %v, want nil", order)
	}
	var cycleErr *CycleError[string, string, string]
	if !errors.As(err, &cycleErr) {
		t.Fatalf("TopologicalSort() error = %T, want *CycleError", err)
	}
	if got, want := cycleErr.Cycle.Nodes, []string{"a", "b", "a"}; !reflect.DeepEqual(got, want) {
		t.Errorf("CycleError nodes = %v, want %v", got, want)
	}
}

func TestTopologicalSortRejectsInconsistentGraph(t *testing.T) {
	order, err := TopologicalSort(inconsistentGraph{})
	if !errors.Is(err, ErrInconsistentGraph) {
		t.Fatalf("TopologicalSort() error = %v, want errors.Is(ErrInconsistentGraph)", err)
	}
	if order != nil {
		t.Errorf("TopologicalSort() order = %v, want nil", order)
	}

	var cycleErr *CycleError[string, string, string]
	if errors.As(err, &cycleErr) {
		t.Errorf("TopologicalSort() error = %T, want contract error, not CycleError", err)
	}
}

func TestStronglyConnectedComponents(t *testing.T) {
	g := buildTestSnapshot(t,
		[]string{"a", "b", "c", "d", "e", "f", "g"},
		[]Edge[string, string, string]{
			{ID: "ab", From: "a", To: "b"},
			{ID: "ba", From: "b", To: "a"},
			{ID: "bc", From: "b", To: "c"},
			{ID: "cd", From: "c", To: "d"},
			{ID: "dc", From: "d", To: "c"},
			{ID: "de", From: "d", To: "e"},
			{ID: "ef", From: "e", To: "f"},
			{ID: "fe", From: "f", To: "e"},
		},
	)

	got := StronglyConnectedComponents(g)
	want := [][]string{{"a", "b"}, {"c", "d"}, {"e", "f"}, {"g"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("StronglyConnectedComponents() = %v, want %v", got, want)
	}
}

func TestAcyclicGraph(t *testing.T) {
	g := buildTestSnapshot(t,
		[]string{"a", "b", "c"},
		[]Edge[string, string, string]{
			{ID: "ab", From: "a", To: "b"},
			{ID: "bc", From: "b", To: "c"},
		},
	)

	if cycle, found := FindCycle(g); found {
		t.Errorf("FindCycle() = %#v, true, want false", cycle)
	}
	if HasCycle(g) {
		t.Error("HasCycle() = true, want false")
	}
}

func buildTestSnapshot(
	t *testing.T,
	ids []string,
	edges []Edge[string, string, string],
) *Snapshot[string, string, string, string] {
	t.Helper()
	nodes := make([]Node[string, string], 0, len(ids))
	for _, id := range ids {
		nodes = append(nodes, Node[string, string]{ID: id, Data: "node-" + id})
	}

	g, err := Build(testSource{nodes: nodes, edges: edges})
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	return g
}
