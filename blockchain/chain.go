package blockchain

import (
	"GoBlockChain/db"
	"GoBlockChain/utils"
	"errors"
	"log"
	"sync"
)

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5
)

type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
}

var b *blockchain
var once sync.Once

var ErrNotFound = errors.New("block not found")

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
			b = &blockchain{Height: 0}
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
	b.CurrentDifficulty = block.Difficulty
}

func (b *blockchain) Blocks() []*Block {
	var blocks []*Block
	hashCursor := b.NewestHash

	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)

		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}

	return blocks
}

func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {

	} else {
		return b.CurrentDifficulty
	}

	return 0
}
