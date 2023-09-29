# Nodes <!-- omit in toc -->

- [Sovereign Relayer](#sovereign-relayer)
- [Dependant Relayer](#dependant-relayer)

This document aims to show a high level diagram of the nodes participating in the Pocket Rollkit Celestia network.

It includes the flow of Requests, Data, Transactions, and Blocks.

## Sovereign Relayer

A Sovereign Relayer is one that:

- Sends Txs to the sequencer
  - specified via `--sequencer-node`
- Runs it's own Pocket full/light client to:
  - read on-chain data
  - listen for on-chain events
  - specify via `--pocket-node` that is either on `localhost` or personally owned domain

_NOTE: The diagram below shows an example where `pocket-node != localhost` but assume it is on a domain owned by the relayer_

```mermaid
---
title: Sovereign Relayer
---
flowchart TB
    a(("Application"))
    subgraph ps["Sequencer"]
        %% Can this be a light client
        pl1[("Pocket Full Node")]
    end
    subgraph r["Relayer (off-chain)"]
        direction TB
        eth[["Ethereum"]]
        gn[["Gnosis"]]
        pg[["Polygon"]]
        etc[["..."]]
    end
    subgraph s["Servicer"]
        pl2[("Pocket Light/Full Node")]
    end
    c{"Celestia DA"}
    a -- Relay Req/Res \n (JSON-RPC endpoint) --> r
    s -. Block & Tx Events \n (Websocket listener).-> r
    r -- Session Dispatch Req/Res \n (JSON-RPC endpoint)--> s
    r -. Txs \n (JSON-RPC endpoint).-> ps
    ps -. Blocks (Commit) .->c
    c -. Blocks (Sync) .-> ps
    c -. Blocks (Sync).-> s
```

## Dependant Relayer

A Dependant Relayer is one that:

- Sends Txs to the sequencer
  - specified via `--sequencer-node`
- Trusts another node (sequencer or other party) to:
  - read on-chain data
  - listen for on-chain events
  - specify via `--pocket-node` that is not on a personally owned domain

_NOTE: The diagram below shows an example where `sequencer-node == pocket-node`_

```mermaid
---
title: Dependant Relayer
---
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
