#!/bin/bash
geth --datadir node2 --networkid 1 --port 30304 --rpc --rpcport "8181" --rpccorsdomain "*" --rpcapi "eth,web3" --bootnodes enode://669c9e350e8d254ee2fc5434b55662badecbcc59c51cda41eb113c1a0722a44654a4cd67d9d400ec93b065392076dc3da43471a706f58b742d434b5315ffb2d8@127.0.0.1:30303
