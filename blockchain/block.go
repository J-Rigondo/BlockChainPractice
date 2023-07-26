package blockchain

import (
	"GoBlockChain/db"
	"GoBlockChain/utils"
	"errors"
	"strings"
	"time"
)

type Block struct {
	Data       string `json:"data"`
	Hash       string `json:"hash"`               // 이전 block의 hash + 현 block의 data
	PrevHash   string `json:"prevHash,omitempty"` //omitempty는 값이 있을 경우만 json에 포함
	Height     int    `json:"height"`
	Difficulty int    `json:"difficulty"`
	Nonce      int    `json:"nonce"`
	Timestamp  int    `json:"timestamp"`
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

func createBlock(data string, prevHash string, height int) *Block {
	block := Block{Data: data, PrevHash: prevHash, Height: height, Difficulty: GetBlockchain().difficulty(), Nonce: 0}

	//payload := block.Data + block.PrevHash + fmt.Sprint(block.Height)
	//block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	block.mine()
	block.persist()

	return &block
}

func FindBlock(hash string) (*Block, error) { //Block 주소를 반환
	blockBytes := db.Block(hash)

	if blockBytes == nil {
		return nil, errors.New("block not found")
	}

	block := &Block{}
	block.restore(blockBytes)

	return block, nil
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)

	for {
		b.Timestamp = int(time.Now().Unix())
		hash := utils.Hash(b)

		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}

	}
}
