# Zomb-AMM

## Decription  
This repository is a added amm module on cosmos-sdk.
I followed B-harvest's lectures of create-amm-module named chaos.

## Important Check point  
I faced a lot of protobuf generating issues during this lectures.  

1. must use [protoc-gen-grpc-gateway@v1.16.0](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.16.0)  
2. must use [protoc-gen-cosmos](https://github.com/regen-network/cosmos-proto)  

## How to start  
I made a start script at root. Just use `sh start.sh`

```
// create tokenA, tokenB liquidity pool
zombd tx amm add-liquidity 10000tokenA,10000tokenB --from zomb-master --keyring-backend test

// check the pool
zombd q amm pairs

// swap tokenA to tokenB
zombd tx amm swap-exact-in 1000tokenA 500tokenB --from zomb-master --keyring-backend test 

// check after balances
zombd q bank balances $(zombd keys show zomb-master -a --keyring-backend test)
```