import "../styles/globals.css";
import type { AppProps } from "next/app";
import { ChainProvider } from "@cosmos-kit/react";
import { ChakraProvider } from "@chakra-ui/react";
import { wallets as keplrWallets } from "@cosmos-kit/keplr";
import { wallets as cosmostationWallets } from "@cosmos-kit/cosmostation";
import { wallets as leapWallets } from "@cosmos-kit/leap";

import { SignerOptions } from "@cosmos-kit/core";
import { chains, assets } from "chain-registry";
import { defaultTheme } from "../config";
import "@interchain-ui/react/styles";

import { chain } from "../config/chain";
import { assetlist } from "../config/assetlist";

function CreateCosmosApp({ Component, pageProps }: AppProps) {
  const signerOptions: SignerOptions = {
    // signingStargate: () => {
    //   return getSigningCosmosClientOptions();
    // }
  };

  return (
    <ChakraProvider theme={defaultTheme}>
      <ChainProvider
        chains={[...chains, chain]}
        assetLists={[...assets, assetlist]}
        wallets={[...keplrWallets, ...cosmostationWallets, ...leapWallets]}
        walletConnectOptions={{
          signClient: {
            projectId: "a8510432ebb71e6948cfd6cde54b70f7",
            relayUrl: "wss://relay.walletconnect.org",
            metadata: {
              name: "Celestia + Cosmos SDK",
              description: "Celestia + Cosmos SDK",
              url: "https://docs.celestia.org/",
              icons: [],
            },
          },
        }}
        endpointOptions={{
          isLazy: true,
        }}
        wrappedWithChakra={true}
        signerOptions={signerOptions}
      >
        <Component {...pageProps} />
      </ChainProvider>
    </ChakraProvider>
  );
}

export default CreateCosmosApp;
