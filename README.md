# Graph Replicator

Small example of a client-server application that handles syncing updates from the client's graph to the server's graph.
The graph being used is an unweighted directed graph.

## Building

`make build`

## Running

Spin up the server:

`./bin/server`

Run the client:

`./bin/client`

### Interesting Notes

When disabling the differential sync support on the on the graph execution times are exponentially slower than with the differential sync support enabled. There are some space-time tradeoffs made to enable differential syncs and they are in no way 100% optimized in the project's current state. Further optimizations to the graph as well as the server can be made but what's here is a good baseline.

### What's Not Done

* The graph implementation is far from optimized

* How the graph is encoded/decoded is far from optimized

* Detecting changes to destination vertices is not optimized

* No persistence (everything is always in memory regardless of how big the graph is)
