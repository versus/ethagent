package main

import "fmt"
import (
	"github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/core/types"
	"log"
	"context"
	"time"
)

func main()  {
    var  quit chan int
	newHead := make(chan *types.Header, 10)
	ctx, _ := context.WithTimeout(context.Background(), 3 * time.Second)

	fmt.Println("Hello World")
	//conn, err := ethclient.Dial("/home/versus/geth/private-chain-data-node1/geth.ipc")
	conn, err := ethclient.Dial("/home/versus/geth/node3/geth.ipc")
	//conn, err := ethclient.Dial("http://127.0.0.1:8080")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	header , err := conn.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Fatalf("Failed get HeaderByNumber: %v", err)
	}

	fmt.Println(header.Number)

	sub,err := conn.SubscribeNewHead(ctx, newHead)
	if err != nil {
		log.Fatalf("Failed : SubscribeNewHead %v", err)
	}

	for {
		select {
		case  <-newHead:
			block := <-newHead
			fmt.Println(block)
		case <-quit:
			fmt.Println("quit")
			return
		}
	}


	sub.Unsubscribe()



}
