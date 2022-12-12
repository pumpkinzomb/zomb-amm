#!/bin/bash
set -eu

PATH=build:$PATH

MONIKER=zomb-node

zombd init $MONIKER --chain-id test-chain-1

zombd keys add zomb-master --keyring-backend test

# Put the generated address in a variable for later use.
MY_VALIDATOR_ADDRESS=$(zombd keys show zomb-master -a --keyring-backend test)

zombd add-genesis-account $MY_VALIDATOR_ADDRESS 100000000stake,100000tokenA,100000tokenB

# Create a gentx.
zombd gentx zomb-master 70000000stake --chain-id test-chain-1 --keyring-backend test

# Add the gentx to the genesis file.
zombd collect-gentxs

zombd start