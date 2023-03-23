package blockchain

import (
	"GoBlockChain/db"
	"GoBlockChain/utils"
	"log"
	"sync"
)

type blockchain struct {
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
}

var b *blockchain
var once sync.Once

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) persist() {
	db.SaveBlockchain(utils.ToBytes(b))

}

func GetBlockchain() *blockchain {
	if b == nil {
		//병렬 실행시 단 한번만 실행하도록
		once.Do(func() {
			b = &blockchain{"", 0}
			checkpoint := db.Checkpoint()
			log.Printf("NewestHash: %s \n Height:%d", b.NewestHash, b.Height)

			if checkpoint == nil {
				b.AddBlock("Genesis Block")

			} else {
				//restore blockchain
				log.Println("restoring...")

				b.restore(checkpoint)

			}
		})
	}
	log.Printf("NewestHash: %s\n Height:%d", b.NewestHash, b.Height)

	return b
}

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height

}
