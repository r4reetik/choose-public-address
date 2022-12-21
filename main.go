package main

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func bruteforce(process *int, key string, keyLen int, mu *sync.RWMutex) {
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

	mu.Lock()
	*process++
	mu.Unlock()

	if address[2:keyLen] == key {
		fmt.Println("Private Key:", hexutil.Encode(privateKeyBytes))
		fmt.Println("Public Addr:", address)
		fmt.Println("------------------------------------------------------------------")
	}
}

func main() {
	mu := &sync.RWMutex{}

	count := 0
	process := 0
	lag := 0
	
	key := string(bytes.ToLower([]byte(os.Args[1])))
	keyLen := len(key) + 2
	
	do := true
	
	fmt.Println("------------------------------------------------------------------")
	for {
		mu.RLock()
		lag = count - process
		fmt.Printf("Count: %v, Process: %v, Lag: %v    \r", count, process, lag)
		mu.RUnlock()

		if lag > 99999 {
			do = false
		} else if lag < 10 {
			do = true
		}

		if do {
			go bruteforce(&process, key, keyLen, mu)
			count++
		}
	}
}
