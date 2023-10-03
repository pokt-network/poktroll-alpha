# Flows <!-- omit in toc -->

## Relay Execution Flow

This flow demonstrations the execution of a `Relay` from the `Application` to the `Miner`. It is a concurrent loop that executes relays in parallel.

`Relay`s are persisted in the `SMST` (Sparse Merkle Sum Tree) to be later included in the `Claim`/`Proof` lifecycle.

The `SessionManager` is a map of `SessionId` -> `SessionWithTree`.

The `SessionWithTree` is a struct that contains the `Session` and the `SMST` of executed `Relay`s.

```mermaid
sequenceDiagram
  title: Relay Execution Flow

  participant Application
  participant Relayer
  participant Servicer as Servicer<br>(Pocket Client)
  participant Miner
  participant SessionManager

  loop Execute Relays (concurrent loop)
    Application ->>+ Relayer: Send signed Relay request
    Relayer ->>+ Servicer: Request current session info
    Servicer ->>- Relayer: Reply with session info
    Relayer ->> Relayer: Execute Relay request
    Relayer ->>- Application: Send signed Relay response
    Relayer -->> Miner: Notify about processed Relay <br> (SignedRequest, SignedResponse)
    Miner ->>+ SessionManager: Request SessionWithTree using <br> Relay.SessionId
    SessionManager ->> SessionManager: Get or Create SessionWithTree <br> for Relay.SessionId
    SessionManager ->>- Miner: Return SMST
    Miner ->> Miner: Insert (SignedRequest, SignedResponse) <br> into SMST
  end
```

## Session Execution Flow

This flow demonstrates the `Claim` & `Proof` lifecycle for `Session`s that are about to finish (i.e. close)

It is a concurrent loop that executes `Claim`/`Proof` submission in parallel across all `Session`s that are closing at the current block height.

```mermaid
sequenceDiagram
  title: Closing Session Flow

  participant Servicer as Servicer<br>(Pocket Client)
  participant Miner
  participant SessionManager

  loop Listen for new blocks
      Servicer -->> SessionManager: Notify of a new committed block
      SessionManager ->> SessionManager: Collect sessions that <br> close at that block height
      SessionManager -->> Miner: Notify Miner with a batch of sessions to close

      loop SessionWithTree map (concurrent loop)
          Miner ->>+ Miner: Get smst_root and close <br> the Session Tree Storestore
          Miner ->>+ Servicer: Submit Tx(MsgClaim)
          note over Miner: Wait for Tx to be included in a block
          Servicer -->> Miner: Tx Event
          Miner ->> Miner: Wait proof wait time blocks<br> (governance parameter)
          Miner ->>+ Servicer: Request latest block hash after waiting
          Servicer ->>- Miner: Reply with latest block hash
          Miner ->> Miner: Re-open Tree Store for proof generation
          Miner ->> Miner: Generate proof with latest block hash
          Miner ->>+ Servicer: Submit Tx(MsgProof)
          note over Miner: Wait for Tx to be included in a block
          Servicer -->> Miner: Tx Event
          Miner ->> Miner: Delete Session Tree
          Miner ->> SessionManager: Delete Session Entry
      end
  end
```
