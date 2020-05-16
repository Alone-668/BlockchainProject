package main

import (
	"BlockchainProject/bolt"
	"log"
)

type BlockchainIterator struct {
	DB *bolt.DB
	curternHashPointer []byte
}

func (bc *BlockChain)NewItertor() *BlockchainIterator {
	it :=BlockchainIterator{
		DB: bc.db,
		curternHashPointer: bc.tail,
	}
	return &it
}

func (it *BlockchainIterator)Next() *Block {
	db := it.DB
	var block Block
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Blockbucket))
		if bucket == nil {
			log.Panic("blockbucket应该存在，请检查数据库")
		}
		blockTem := bucket.Get(it.curternHashPointer)
		block = Deserialize(blockTem)

		it.curternHashPointer = block.PrevHash
		return nil

	})
	return &block

}