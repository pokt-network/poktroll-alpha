This document describes claim and proof generation for sessions that are closing.

It is a concurrent loop that executes the claim/proof submission in parallel across all sessions that are closing at the current block height.

```mermaid
sequenceDiagram
  title: Closing sessions handling flow
  participant Servicer as Servicer<br>(Pocket Client)
  participant Miner
  participant SessionManager

  loop Listen to new blocks
      Servicer -->> SessionManager: Notify of a new committed block
      SessionManager ->> SessionManager: Collect sessions that close at that block height
      SessionManager -->> Miner: Notify miner with a batch of sessions to close

      loop SessionWithTree map (concurrent loop)
          Miner ->>+ Miner: Get smst_root and close the tree store
          Miner ->>+ Servicer: Submit Tx(MsgClaim)
          note over Miner: Wait for Tx to be included in a block
          Miner ->> Miner: Wait proof wait time blocks<br> (governance parameter)
          Miner ->>+ Servicer: Request latest block hash after waiting
          Servicer ->>- Miner: Reply with latest block hash
          Miner ->> Miner: Re-open tree store for proof generation
          Miner ->> Miner: Generate proof with latest block hash
          Miner ->>+ Servicer: Submit proof
          note over Miner: Wait for Tx to be included in a block
          Miner ->> Miner: Delete session's tree
          Miner ->> SessionManager: Delete session's entry
      end
  end
```