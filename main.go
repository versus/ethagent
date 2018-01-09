package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"os"

	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/BurntSushi/toml"
	"flag"
)

type Config struct {
    Endpoint string
    Hostname string
    IPCPath	 string `toml:"ipcpath"`
    Backup 	 bool
}

func main() {
	flagConfigFile := flag.String("c", "./config.toml", "config: path to config file")
	flag.Parse()
	sig := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	newHead := make(chan *types.Header, 10)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()



	log.Println("ethagent v 0.0.1")
	var conf Config
	if _, err := toml.DecodeFile(*flagConfigFile, &conf); err != nil {
		log.Fatalln("Error parse config.toml", err.Error())
	}

	conn, err := ethclient.Dial(conf.IPCPath)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	header, err := conn.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Fatalf("Failed get HeaderByNumber: %v", err)
	}

	fmt.Println(header.Number)

	sub, err := conn.SubscribeNewHead(ctx, newHead)
	if err != nil {
		log.Fatalf("Failed : SubscribeNewHead %v", err)
	}

	go func() {
		for {
			select {
			case <-newHead:
				block := <-newHead
				fmt.Println(block)
			case <-sig:
				s := <-sig
				fmt.Println(s)
				done <- true
				return
			}
		}
	}()

	<-done
	sub.Unsubscribe()

}
