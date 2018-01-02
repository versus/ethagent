mkdir ./node{1,2}
geth --datadir node1 account new
geth --datadir node1 init genesis.json
geth --datadir node2 account new
geth --datadir node2 init genesis.json
