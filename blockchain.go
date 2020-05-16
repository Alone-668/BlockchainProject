package main

import (
	"BlockchainProject/bolt"
	"log"
)

type BlockChain struct {
	//定一个区块链数组
	//blocks []*Block
	db *bolt.DB
	tail []byte
}

const BlockchainDB = "BlockChain.db"
const Blockbucket = "blockBucket"
//5. 定义一个区块链
func NewBlockChain() *BlockChain {
	var lastHash []byte
	//创建一个创世块，并作为第一个区块添加到区块链中
	db, err := bolt.Open(BlockchainDB, 0600, nil)
	if err != nil {
		log.Panic("打开数据库失败！")
	}
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Blockbucket))
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(Blockbucket))
			if err != nil {
				log.Panic("bucket创建失败")
			}
			genesisBlock := GenesisBlock()
			bucket.Put([]byte(genesisBlock.Hash),genesisBlock.Serialize())
			//用来存储最后一个区块的哈希值，方便从数据库中遍历出区块信息
			bucket.Put([]byte("LastHashKey"),genesisBlock.Hash)
			lastHash = genesisBlock.Hash
		}else {
			lastHash = bucket.Get([]byte("LastHashKey"))
		}
		return nil
	})

	return &BlockChain{
		db: db,
		tail: lastHash,
	}
}
func GenesisBlock() *Block {
	return NewBlock("创世块", []byte{})
}

//5. 添加区块
func (bc *BlockChain) AddBlock(data string) {

	db := bc.db
	lastHash := bc.tail

	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Blockbucket))
		if bucket == nil {
			log.Panic("blockBucket不应该为空，bolt数据库出现问题")
		}
		block := NewBlock(data, lastHash)
		bucket.Put(block.Hash,block.Serialize())
		bucket.Put([]byte("LastHashKey"),block.Hash)
		bc.tail = block.Hash

		return nil
	})


}

