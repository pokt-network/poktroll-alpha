# Pocket Network RollKit Alpha <!-- omit in toc -->

**IMPORTANT: DO NOT GO TO PRODUCTION WITH THIS REPOSITORY**

This is an alpha version of trying to build Pocket on top of RollKit.

- [Reference TUtorial](#reference-tutorial)
- [Getting Started](#getting-started)

## Reference TUtorial

Following the tutorial here: https://rollkit.dev/tutorials

## Getting Started

Start up a celestia localnet:

```bash
    make celestia_localnet
```

Export the DA token:

```bash
    export CELESTIA_NODE_AUTH_TOKEN=$(make celestia_localnet_auth_token)
```

Try retrieving the balance:

```bash
    make celestia_localnet_balance_check
```

