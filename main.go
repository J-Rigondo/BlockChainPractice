package main

import (
	"GoBlockChain/blockchain"
)

func main() {
	chain := blockchain.GetBlockchain()
	chain.AddBlock("jun")
	chain.AddBlock("chan")
	chain.ListBlocks()

}
