#!/bin/bash

geth --datadir ./node1/  --rpc --rpcport "8080" --rpccorsdomain "*" --rpcapi "db,eth,net,web3,personal,miner" --ws --wsorigins "*" --wsapi "db,eth,net,web3,personal,miner"
