package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
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
	//b. 交易数组
	//Data []byte
	Transactions []*Transaction
}

func (block *Block) Serialize() []byte {
	var bufer bytes.Buffer
	encoder := gob.NewEncoder(&bufer)
	err := encoder.Encode(&block)
	if err != nil {
		log.Panic("编码失败")
	}
	return bufer.Bytes()
}
//反序列化
func Deserialize(data []byte) Block {

	decoder := gob.NewDecoder(bytes.NewReader(data))

	var block Block
	//2. 使用解码器进行解码
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic("解码出错!")
	}

	return block
}

func NewBlock(txs []*Transaction, prevBlockHash []byte) *Block {
	block := Block{
		Version:    00,
		PrevHash:   prevBlockHash,
		MerkelRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 00,
		Nonce:      00,
		Hash:       []byte{},
		Transactions: txs,
	}
	//block.setHash()
	block.MerkelRoot = block.MakeMerkelRoot()
	pow := Newproofofwork(&block)
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
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

/*func (block *Block) setHash() {
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

}*/
//模拟梅克尔根，只是对交易的数据做简单的拼接，而不做二叉树处理！
//func (block *Block) MakeMerkelRoot() []byte {
func (block *Block)MakeMerkelRoot() []byte {
	var info []byte
	for _,tx:= range block.Transactions{
		//拼接整体交易的ID（交易结构体Hash）
		info = append(info, tx.TXID...)
	}
	hash := sha256.Sum256(info)
	return hash[:]

}