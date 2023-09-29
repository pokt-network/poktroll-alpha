This document describes the flow of relay execution from the application to the miner. It is a concurrent loop that executes relays in parallel.

Relays are persisted into the SMST for later claim/proof generation. The SMST is a map of SessionId -> SessionWithTree.

The SessionWithTree is a struct that contains the Session and the SMST of the executed relays.

```mermaid
sequenceDiagram
  title: Relay execution flow
  participant Application
  participant Relayer
  participant Servicer as Servicer<br>(Pocket Client)
  participant Miner
  participant SessionManager

  loop Execute Relays (concurrent loop)
    Application ->>+ Relayer: Send signed relay request
    Relayer ->>+ Servicer: Request current session info
    Servicer ->>- Relayer: Reply with session info
    Relayer ->> Relayer: Execute relay request
    Relayer ->>- Application: Send signed relay response
    Relayer -->> Miner: Notify about processed relay
    Miner ->>+ SessionManager: Request SessionWithTree using<br>relay's sessionId
    SessionManager ->> SessionManager: Create SessionWithTree at<br>closing-block entry if not exists
    SessionManager ->>- Miner: Return SMST
    Miner ->> Miner: Insert into SMST
  end
```