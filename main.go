package main

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func bruteforce(process *uint, key string, keyLen int) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	address := string(bytes.ToLower([]byte(crypto.PubkeyToAddress(*publicKeyECDSA).Hex())))

	*process++

	if address[2:keyLen] == key {
		fmt.Println("Public Addr:", address)
		fmt.Println("Private Key:", hexutil.Encode(privateKeyBytes))
		fmt.Println("------------------------------------------------------------------")
	}
}

func main() {
	var count uint = 0
	var process uint = 0
	var lag uint = 0
	var key string = os.Args[1]
	var keyLen int = len(key) + 2
	do := true
	for {
		lag = count - process
		fmt.Printf("Count: %v, Process: %v, Lag: %v\r", count, process, lag)

		if lag > 9999 {
			do = false
		} else if lag < 10 {
			do = true
		}

		if do {
			go bruteforce(&process, key, keyLen)
			count++
		}
	}
}
