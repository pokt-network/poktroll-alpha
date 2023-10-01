# Nodes <!-- omit in toc -->

- [Dependant Relayer](#dependant-relayer)
- [Sovereign Relayer](#sovereign-relayer)

This document aims to show a high level diagram of the nodes participating in the Pocket Rollkit Celestia network.

It includes the flow of Requests, Data, Transactions, and Blocks.

## Dependant Relayer

The diagram below shows the absolute base case where there is:

1. 1 Pocket Rollup Node
2. The Rollup Node is also the Centralized Sequencer
3. The Centralized Sequencer is also the Proxy's (i.e. Relayer/Miner) source of data and events

A Dependant Relayer is one that:

- Sends Txs to the sequencer
  - specified via `--sequencer-node`
- Trusts another node (i.e. spefieid via `--pocket-node`) to:
  - read on-chain data
  - listen for on-chain events

```mermaid
flowchart TB
    a(("Application"))
    subgraph p["Pocket Node"]
        direction LR
        rs([Role 1 - Sequencer])
        rv([Role 2 - Servicer])
        pl1[("Pocket Full Node")]
    end
    subgraph r["Relayer (off-chain)"]
        direction TB
        eth[["Ethereum"]]
        gn[["Gnosis"]]
        pg[["Polygon"]]
        etc[["..."]]
    end
    c{"Celestia DA"}
    a -- RPC Relay Req/Res \n (JSON-RPC endpoint) --> r
    p -. Block & Tx Events \n (Websocket listener).-> r
    r -- Session Dispatch Req/Res \n (JSON-RPC endpoint)--> p
    r -. Txs \n (JSON-RPC endpoint).-> p
    p -. Blocks (Commit) .->c
    c -. Blocks (Sync) .-> p
```

## Sovereign Relayer

The diagram below shows the Celestia DA, Rollup Nodes in the network, the sequencer as well as a Sovereign Relayer that maintains its own Pocket Rollup Node.

A Sovereign Relayer is one that:

- Sends Txs to the sequencer
  - specified via `--sequencer-node`
- Runs it's own Pocket Full Node (specified via `--pocket-node`) to:
  - read on-chain data
  - listen for on-chain events

```mermaid
flowchart TB
    a(("Application"))
    subgraph prln["Rollup Nodes"]
        pfn1[("Pocket Full Node")]
        pfn2[("Pocket Full Node")]
        pfn3[("Pocket Full Node")]
        pfn1 <-. gossip \n (Txs & Blocks) .-> pfn2
        pfn2 <-. gossip \n (Txs & Blocks) .-> pfn3
        pfn3 <-. gossip \n (Txs & Blocks) .-> pfn1
    end
    subgraph ps["Sequencer"]
        pl1[("Pocket Full Node")]
    end
    subgraph r["Proxy (off-chain Relayer & Miner)"]
        direction TB
        eth[["Ethereum"]]
        gn[["Gnosis"]]
        pg[["Polygon"]]
        etc[["..."]]
    end
    subgraph s["Servicer (Rollup Node maintainedby Proxy Operator) "]
        pl2[("Pocket Full Node")]
    end
    c{"Celestia DA"}
    a -- Relay Req/Res \n (JSON-RPC endpoint) --> r
    s -. Block & Tx Events \n (Websocket listener).-> r
    r -- Session Dispatch Req/Res \n (JSON-RPC endpoint)--> s
    r -. Txs \n (JSON-RPC endpoint).-> prln
    r -. Txs \n (JSON-RPC endpoint).-> ps
    prln <-. gossip \n (Txs & Blocks) .-> ps
    ps -. Blocks\n(Commit) .->c
    c -. Blocks\n(Sync) .-> ps
    c -. Blocks\n(Sync).-> s
```
