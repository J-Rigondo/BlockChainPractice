package main

import (
	"GoBlockChain/blockchain"
	"log"
)

func main() {
	blockchain.GetBlockchain()

	block, _ := blockchain.FindBlock("81f2ced897805e5539e750784e8d12bff104712be9bf8845ce52e006b0f3252e")
	log.Println(block)

}
