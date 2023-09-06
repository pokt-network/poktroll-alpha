# POKTRoll <!-- omit in toc -->

**IMPORTANT: DO NOT GO TO PRODUCTION WITH THIS REPOSITORY**

This is an alpha version of trying to build Pocket on top of RollKit.

**poktroll** is a rollup built using Cosmos SDK and Tendermint and created with [Ignite CLI](https://ignite.com/cli).

- [LocalNet](#localnet)
- [Testnet](#testnet)
- [AUTOGENERATED README BELOW](#autogenerated-readme-below)
- [Get started](#get-started)
  - [Configure](#configure)
  - [Web Frontend](#web-frontend)
- [Release](#release)
  - [Install](#install)
- [Learn more](#learn-more)

## LocalNet

Following the tutorial here: https://rollkit.dev/tutorials

Start up a celestia localnet:

```bash
    make celestia_localnet
```

Try retrieving your Celestia token balance:

```bash
    make celestia_localnet_balance_check
```

Start up a poktroll localnet:

```bash
    make poktroll_local_start
```

## Testnet

Start up a celestia light node:

```bash
    make celestia_light_client_start
```

Start up a poktroll testnet:

```bash
    make poktroll_testnet_start
```

Try retrieving your Celestia testnet balance:

```bash
    make celestia_testnet_balance_check
```

## Poktroll commands

Send tokens on poktroll:

```bash
    make poktroll_send
```

Get token balances after/before transfer:

```bash
    make poktroll_balance
```

Get session data:

```bash
    make poktroll_get_session
```

## AUTOGENERATED README BELOW

## Get started

```bash
ignite chain serve
```

`serve` command installs dependencies, builds, initializes, and starts your blockchain in development.

### Configure

Your blockchain in development can be configured with `config.yml`. To learn more, see the [Ignite CLI docs](https://docs.ignite.com).

### Web Frontend

Ignite CLI has scaffolded a Vue.js-based web app in the `vue` directory. Run the following commands to install dependencies and start the app:

```bash
cd vue
npm install
npm run serve
```

The frontend app is built using the `@starport/vue` and `@starport/vuex` packages. For details, see the [monorepo for Ignite front-end development](https://github.com/ignite/web).

## Release

To release a new version of your blockchain, create and push a new tag with `v` prefix. A new draft release with the configured targets will be created.

```bash
git tag v0.1
git push origin v0.1
```

After a draft release is created, make your final changes from the release page and publish it.

### Install

To install the latest version of your blockchain node's binary, execute the following command on your machine:

```bash
curl https://get.ignite.com/username/poktroll@latest! | sudo bash
```

`username/poktroll` should match the `username` and `repo_name` of the Github repository to which the source code was pushed. Learn more about [the install process](https://github.com/allinbits/starport-installer).

## Learn more

- [Ignite CLI](https://ignite.com/cli)
- [Tutorials](https://docs.ignite.com/guide)
- [Ignite CLI docs](https://docs.ignite.com)
- [Cosmos SDK docs](https://docs.cosmos.network)
- [Developer Chat](https://discord.gg/ignite)
