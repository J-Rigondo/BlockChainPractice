package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type block struct {
	data     string
	hash     string // 이전 block의 hash + 현 block의 data
	prevHash string
}

type blockchain struct {
	blocks []*block //private
}

var b *blockchain
var once sync.Once

func GetBlockchain() *blockchain {
	if b == nil {
		//병렬 실행시 단 한번만 실행하도록
		once.Do(func() {
			b = &blockchain{}
			b.AddBlock("Genesis Block")
		})
	}

	return b
}

func (b *blockchain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

func createBlock(data string) *block {
	newBlock := block{data, "", getLastHash()}
	sum256 := sha256.Sum256([]byte(data + newBlock.prevHash))
	hexHash := fmt.Sprintf("%x", sum256)
	newBlock.hash = hexHash

	return &newBlock

}

func getLastHash() string {
	if len(b.blocks) > 0 {
		return b.blocks[len(b.blocks)-1].hash
	}

	return ""
}

func (b *blockchain) ListBlocks() {
	for _, a := range b.blocks {
		fmt.Printf("data: %s\n", a.data)
		fmt.Printf("hash: %s\n", a.hash)
		fmt.Printf("prevHash: %s\n", a.prevHash)
	}
}
