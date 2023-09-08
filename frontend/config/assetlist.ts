export const assetlist = {
  $schema: "../../assetlist.schema.json",
  chain_name: "poktroll",
  assets: [
    {
      description: "",
      denom_units: [
        {
          denom: "ustake",
          exponent: 0,
        },
        {
          denom: "stake",
          exponent: 6,
        },
      ],
      base: "stake",
      name: "Pocket Network Rollup",
      display: "stake",
      symbol: "STAKE",
      logo_URIs: {
        png: "https://www.pokt.network/wp-content/uploads/2023/04/logo-small.png",
      },
    },
  ],
};
