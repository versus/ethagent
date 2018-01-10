package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"os"

	"os/signal"
	"syscall"

	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Config struct {
	Endpoint string `json:"-"`
	Hostname string `json:"hostname"`
	IPCPath  string `toml:"ipcpath" json:"-"`
	State    string `json:"state"`
	Block    big.Int `json:"block"`
}

var conf Config

func sendNewBlock() {
	jsonStr, err := json.Marshal(&conf)
	if err != nil {
		log.Fatalf("Failed to convert Conf to json", err.Error())
	}
	var cancel context.CancelFunc
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequest("POST", conf.Endpoint, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body ",string(body))
}

func main() {
	flagConfigFile := flag.String("c", "./config.toml", "config: path to config file")
	flag.Parse()
	sig := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	newHead := make(chan *types.Header, 10)


	log.Println("ethagent v 0.0.1")

	if _, err := toml.DecodeFile(*flagConfigFile, &conf); err != nil {
		log.Fatalln("Error parse config.toml", err.Error())
	}

	conn, err := ethclient.Dial(conf.IPCPath)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
    ctx := context.Background()
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
				conf.Block = *block.Number
				sendNewBlock()
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
