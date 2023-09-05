#!/bin/sh
echo "AUTH_TOKEN is: $AUTH_TOKEN"

# set variables for the chain
VALIDATOR_NAME=validator1
CHAIN_ID=poktroll
KEY_NAME=poktroll-key
KEY_2_NAME=poktroll-key-2
CHAINFLAG="--chain-id ${CHAIN_ID}"
TOKEN_AMOUNT="10000000000000000000000000stake"
STAKING_AMOUNT="1000000000stake"
AUTH_TOKEN=$AUTH_TOKEN

echo "~~~~~"
echo $AUTH_TOKEN
echo "~~~~~"

# create a random Namespace ID for your rollup to post blocks to
NAMESPACE_ID=$(openssl rand -hex 10)

# query the DA Layer start height, in this case we are querying
# our local devnet at port 26657, the RPC. The RPC endpoint is
# to allow users to interact with Celestia's nodes by querying
# the node's state and broadcasting transactions on the Celestia
# network. The default port is 26657.
DA_BLOCK_HEIGHT=$(curl http://0.0.0.0:26657/block | jq -r '.result.block.header.height')

# rollkit logo
cat <<'EOF'

                 :=+++=.
              -++-    .-++:
          .=+=.           :++-.
       -++-                  .=+=: .
   .=+=:                        -%@@@*
  +%-                       .=#@@@@@@*
    -++-                 -*%@@@@@@%+:
       .=*=.         .=#@@@@@@@%=.
      -++-.-++:    =*#@@@@@%+:.-++-=-
  .=+=.       :=+=.-: @@#=.   .-*@@@@%
  =*=:           .-==+-    :+#@@@@@@%-
     :++-               -*@@@@@@@#=:
        =%+=.       .=#@@@@@@@#%:
     -++:   -++-   *+=@@@@%+:   =#*##-
  =*=.         :=+=---@*=.   .=*@@@@@%
  .-+=:            :-:    :+%@@@@@@%+.
      :=+-             -*@@@@@@@#=.
         .=+=:     .=#@@@@@@%*-
             -++-  *=.@@@#+:
                .====+*-.

   ______         _  _  _     _  _
   | ___ \       | || || |   (_)| |
   | |_/ /  ___  | || || | __ _ | |_
   |    /  / _ \ | || || |/ /| || __|
   | |\ \ | (_) || || ||   < | || |_
   \_| \_| \___/ |_||_||_|\_\|_| \__|
EOF

# echo variables for the chain
echo -e "\n\n\n\n\n Your NAMESPACE_ID is $NAMESPACE_ID \n\n Your DA_BLOCK_HEIGHT is $DA_BLOCK_HEIGHT \n\n\n\n\n"

# build the poktroll chain with Rollkit
ignite chain build

# reset any existing genesis/chain data
poktrolld tendermint unsafe-reset-all

# initialize the validator with the chain ID you set
poktrolld init $VALIDATOR_NAME --chain-id $CHAIN_ID

# add keys for key 1 and key 2 to keyring-backend test
poktrolld keys add $KEY_NAME --keyring-backend test
poktrolld keys add $KEY_2_NAME --keyring-backend test

# add these as genesis accounts
poktrolld add-genesis-account $KEY_NAME $TOKEN_AMOUNT --keyring-backend test
poktrolld add-genesis-account $KEY_2_NAME $TOKEN_AMOUNT --keyring-backend test

# set the staking amounts in the genesis transaction
poktrolld gentx $KEY_NAME $STAKING_AMOUNT --chain-id $CHAIN_ID --keyring-backend test

# collect genesis transactions
poktrolld collect-gentxs

# start the chain
echo "========="
echo $AUTH_TOKEN
poktrolld start --rollkit.aggregator true --rollkit.da_layer celestia --rollkit.da_config='{"base_url":"http://localhost:26658","timeout":60000000000,"fee":600000,"gas_limit":6000000,"auth_token":"'$AUTH_TOKEN'"}' --rollkit.namespace_id $NAMESPACE_ID --rollkit.da_start_height $DA_BLOCK_HEIGHT --rpc.laddr tcp://127.0.0.1:36657 --p2p.laddr "0.0.0.0:36656"
echo "========="

# uncomment the next command if you are using lazy aggregation
# poktrolld start --rollkit.aggregator true --rollkit.da_layer celestia --rollkit.da_config='{"base_url":"http://localhost:26658","timeout":60000000000,"fee":600000,"gas_limit":6000000,"auth_token":"'$AUTH_TOKEN'"}' --rollkit.namespace_id $NAMESPACE_ID --rollkit.da_start_height $DA_BLOCK_HEIGHT --rollkit.lazy_aggregator
