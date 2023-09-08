export const chain = {
  $schema: "../../chain.schema.json",
  chain_name: "poktroll",
  chain_id: "poktroll",
  pretty_name: "pocket network rollup",
  status: "live",
  network_type: "testnet",
  bech32_prefix: "poktroll",
  daemon_name: "poktrolld",
  node_home: "$HOME/.poktroll",
  key_algos: ["secp256k1"],
  slip44: 118,
  fees: {
    fee_tokens: [
      {
        denom: "ustake",
        fixed_min_gas_price: 0,
      },
    ],
  },
  apis: {
    rpc: [
      {
        address: "http://localhost:26657",
        provider: "JCS",
      },
    ],
    rest: [
      {
        address: "http://localhost:1317",
        provider: "JCS",
      },
    ],
  },
  beta: true,
};
