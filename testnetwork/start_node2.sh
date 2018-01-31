#!/bin/bash
geth --datadir node2 --networkid 15 --port 30304 --rpc --rpcport "8181" --rpccorsdomain "*" --rpcapi "eth,web3" --bootnodes enode://52f6c5f7e69ebd494361f064aac3166e75e0875c1de365767dd889fdd2d87045780ea1301024e0e732981c9bfe338e5bbc1f2d8ff4721dce61955c103863d300@127.0.0.1:30303

