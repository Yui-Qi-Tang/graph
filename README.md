# graph

Version: v0.1.0 (unreleased)

`graph` is a small, domain-neutral Go package for immutable directed graph
snapshots and structural graph algorithms.

## v0.1 scope

- Generic node and independently identified edge records.
- Validated, immutable topology snapshots with stable iteration order.
- Incoming and outgoing edge inspection.
- BFS, DFS, reachability, and shortest unweighted path witnesses.
- Cycle detection with edge-complete witnesses.
- Deterministic topological sorting.
- Strongly connected components.

Node and edge data are transported without interpretation. Snapshot topology is
copied during construction, while caller-defined data has shallow-copy semantics.

The package intentionally contains no domain rules, persistence, revision model,
SAT encoding, plugin system, or mutation API.

## Example

```go
package main

import (
	"fmt"
	"iter"
	"slices"

	graph "yuki.tang.github.com"
)

type source struct {
	nodes []graph.Node[string, struct{}]
	edges []graph.Edge[string, string, struct{}]
}

func (s source) Nodes() iter.Seq[graph.Node[string, struct{}]] {
	return slices.Values(s.nodes)
}

func (s source) Edges() iter.Seq[graph.Edge[string, string, struct{}]] {
	return slices.Values(s.edges)
}

func main() {
	snapshot, err := graph.Build(source{
		nodes: []graph.Node[string, struct{}]{
			{ID: "fetch"}, {ID: "build"}, {ID: "deploy"},
		},
		edges: []graph.Edge[string, string, struct{}]{
			{ID: "fetch-build", From: "fetch", To: "build"},
			{ID: "build-deploy", From: "build", To: "deploy"},
		},
	})
	if err != nil {
		panic(err)
	}

	order, err := graph.TopologicalSort(snapshot)
	if err != nil {
		panic(err)
	}
	fmt.Println(order) // [fetch build deploy]
}
```

## Verify

```sh
go test -count=1 ./...
go vet ./...
go build ./...
```
