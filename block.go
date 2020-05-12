package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"time"
)

type Block struct {
	//1.版本号
	Version uint64
	//2. 前区块哈希
	PrevHash []byte
	//3. Merkel根（梅克尔根，这就是一个哈希值，我们先不管，我们后面v4再介绍）
	MerkelRoot []byte
	//4. 时间戳
	TimeStamp uint64
	//5. 难度值
	Difficulty uint64
	//6. 随机数，也就是挖矿要找的数据
	Nonce uint64

	//a. 当前区块哈希,正常比特币区块中没有当前区块的哈希，我们为了是方便做了简化！
	Hash []byte
	//b. 数据
	Data []byte
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		Version:    00,
		PrevHash:   prevBlockHash,
		MerkelRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 00,
		Nonce:      00,
		Hash:       []byte{},
		Data:       []byte(data),
	}
	block.setHash()
	return &block

}
func Uint64toByte(num uint64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()

}
func (block *Block) setHash() {
	temp := [][]byte{
		Uint64toByte(block.Version),
		block.PrevHash,
		block.MerkelRoot,
		Uint64toByte(block.TimeStamp),
		Uint64toByte(block.Difficulty),
		Uint64toByte(block.Nonce),
		block.Data,
	}
	blockInfo := bytes.Join(temp, []byte{})
	myhash := sha256.Sum256(blockInfo)

	block.Hash = myhash[:]

}
