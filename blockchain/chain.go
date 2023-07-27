package blockchain

import (
	"GoBlockChain/db"
	"GoBlockChain/utils"
	"errors"
	"fmt"
	"log"
	"sync"
)

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5
	blockInterval      int = 2
	allowedRange       int = 2
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
	b.persist()
}

func (b *blockchain) Blocks() []*Block {
	var blocks []*Block
	hashCursor := b.NewestHash
	// todo b는 주소지만 byte로 변환하려고하면 (*b)이렇게 자동으로 해석하나보다
	log.Println("===============================")
	log.Println(b)
	fmt.Printf("type:  %T", b)

	bytes := utils.ToBytes(b)
	log.Println(bytes)
	temp := &blockchain{}
	utils.FromBytes(temp, bytes)
	log.Println(temp)

	log.Println("===============================")

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

func (b *blockchain) recalculateDifficulty() int {
	allBlocks := b.Blocks()
	newestBlock := allBlocks[0]
	lastRecalculatedBlock := allBlocks[difficultyInterval-1]
	actualTime := (newestBlock.Timestamp / 60) - (lastRecalculatedBlock.Timestamp / 60)
	expectedTime := difficultyInterval * blockInterval
	if actualTime <= (expectedTime - allowedRange) {
		return b.CurrentDifficulty + 1
	} else if actualTime >= (expectedTime + allowedRange) {
		return b.CurrentDifficulty - 1
	}
	return b.CurrentDifficulty
}

func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		return b.recalculateDifficulty()

	} else {
		return b.CurrentDifficulty
	}

	return 0
}
