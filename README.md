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

## Verify

```sh
go test -count=1 ./...
go vet ./...
go build ./...
```
