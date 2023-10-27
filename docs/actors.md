# Actors <!-- omit in toc -->

- [Actor Module Interaction](#actor-module-interaction)
- [Relayer Architecture](#relayer-architecture)

This document aims to show a high level diagram of Pocket Network actors and the interaction between them.

## Actor Module Interaction

This diagram gives an overview of the interaction between the various on/off chain actors.

```mermaid
---
title: Actor/module interaction
---

flowchart

subgraph on["On-chain (modules)"]
    subgraph s[Servicer]
        s_s[Stake]
        s_c[Claims]
        s_p[Proofs]
    end

    subgraph se[Session]
        se_get[Get session]
    end

    subgraph aon[Application]
        a_s[Stake]
        a_a[Portal Refs]
        a_pk[Signing account]
    end

    subgraph pon[Portal]
        p_s[Stake]
        p_pk[Signing account]
    end
end

subgraph off["Off-chain"]
    subgraph r["Relayer (servicer)"]
        r_proxy[Proxies relays]
        r_claim[Submits claims]
        r_proof[Submits proofs]
        r_relay[Signs relay responses]
    end

    subgraph poff[Portal]
    end

    subgraph aoff[Application]
    end

    subgraph sc[Signing Client]
        sign[Signs relay requests]
    end

    subgraph eu[End user]
        eu_r[Requests relays]
    end
end


poff --> pon
poff --> sc
poff --> se
aoff --> sc
aoff --> se
aoff --> aon

sc --> r
r --> s
r --> se

eu --> poff
eu --> aoff
```

## Relayer Architecture

This diagram gives an overview of the core business logic of a Relayer.

```mermaid
---
title: Relayer Architecture
---

flowchart

subgraph client["Signing client (App/Gateway)"]
    signs["Signs & sends relays\n (RPC Requests)"]
    c_res[Receives RPC Responses]
end

subgraph rel[Relayer]
    subgraph proxy[Proxy]
        val_ses["Session Validation (e.g. matches App)"]
    end

    subgraph ses_man[Session Manager]
    end

    subgraph miner[Miner]
        dif[Hashes relays & checks their difficulty]
    end
end

subgraph target["Relay target"]
    req[Request handling]
    res[Response origin]
end

subgraph seq["Sequencer / Validator"]
    subgraph comet[Cometbft node]
        comet_rpc[JSON-RPC endpoint]
        comet_grpc[gRPC endpoint]
    end

    subgraph app[PoktRoll Rollup Node]

        subgraph servicer[Servicer module]
            subgraph ser_claim[Claims]
                %% ser_val_ser[Validates signer matches servicer addr]
            end
            subgraph ser_proof[Proofs]
                %% ser_val_proof[Validates proof]
            end
        end

        subgraph session[Session module]
            ses_get[GetSession]
        end
    end
end

client --> proxy
proxy --"update session SMST"--> ses_man
proxy --"completed relays"--> miner
proxy --"execute relays"--> target
proxy --"get session"--> comet_grpc

ses_man -. "websocket subscription to blocks" .-> comet_rpc
ses_man --> miner

miner --"claim / proof txs"--> comet_grpc

comet --"ABCI (Query)"--> session
comet --"ABCI (DeliverTx)"--> servicer
```
