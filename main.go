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

	"sync"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Config struct {
	Server    string     `toml:"server" json:"-"`
	IPCPath   string     `toml:"ipcpath" json:"-"`
	State     string     `json:"state"`
	Block     big.Int    `json:"block"`
	Token     string     `json:"token"`
	AccessKey string     `toml:"key" json:"-"`
	Mutex     sync.Mutex `json:"-"`
}

var (
	conf         Config
	addrNewblock string
	addrAuth     string
)

const (
	Version           = "v0.0.1"
	Author            = " by Valentyn Nastenko [versus.dev@gmail.com]"
	endpoint_newblock = "/api/v1/newblock"
	endpoint_auth     = "/api/v1/auth"
)

func sendNewBlock(number big.Int) {
	//TODO: send new block to multihost

	var cancel context.CancelFunc
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conf.Mutex.Lock()
	conf.Block = number
	jsonStr, err := json.Marshal(&conf)
	if err != nil {
		log.Fatalf("Failed to convert Conf to json", err.Error())
	}
	conf.Mutex.Unlock()

	req, err := http.NewRequest("POST", addrNewblock, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Println("error request to endpoint ", addrNewblock, err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("error client do ", err.Error())
		return
	}
	resp.Body.Close()

	if *flagDebug {
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("response Body ", string(body))
	}
	//TODO: parse and use last block from server
	resp.Body.Close()
}

func main() {
	// TODO: add prometheus metrics https://stackoverflow.com/questions/37611754/how-to-push-metrics-to-prometheus-using-client-golang
	flagConfigFile := flag.String("c", "./config.toml", "config: path to config file")
	gnrToken := flag.Bool("gentoken", false, "config: generate token for agents")
	flagDebug := flag.Bool("vvv", false, "runtime: output debug messages ")
	flag.Parse()

	if *gnrToken {
		fmt.Println("Token is ", GenToken(16))
		os.Exit(0)
	}

	sig := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	newHead := make(chan *types.Header, 10)

	log.Println("ethagent ", Version, Author)

	if _, err := toml.DecodeFile(*flagConfigFile, &conf); err != nil {
		log.Fatalln("Error parse config.toml", err.Error())
	}

	addrNewblock = fmt.Sprintf("%s%s", conf.Server, endpoint_newblock)
	addrAuth = fmt.Sprintf("%s%s", conf.Server, endpoint_auth)
	if *flagDebug {
		log.Println("addrNewblock is ", addrNewblock)
		log.Println("addrAuth is ", addrAuth)
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
	auth()
	sub, err := conn.SubscribeNewHead(ctx, newHead)
	if err != nil {
		log.Fatalf("Failed : SubscribeNewHead %v", err)
	}

	go func() {
		for {
			select {
			case <-newHead:
				block := <-newHead
				go sendNewBlock(*block.Number)
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
	os.Exit(0)

}
