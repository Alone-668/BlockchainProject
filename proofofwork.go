package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func Newproofofwork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		block: block,
	}
	temTarget := "0000f00000000000000000000000000000000000000000000000000000000000"
	temInt := big.Int{}
	temInt.SetString(temTarget, 16)
	pow.target = &temInt
	return &pow

}

func (pow *ProofOfWork) Run() ([]byte, uint64) {
	var nonce uint64
	block := pow.block
	var hash [32]byte

	for {
		tem := [][]byte{
			Uint64toByte(block.Version),
			block.PrevHash,
			block.MerkelRoot,
			Uint64toByte(block.Difficulty),
			Uint64toByte(block.TimeStamp),
			Uint64toByte(nonce),
			//block.Data,
		}
		blockInfo := bytes.Join(tem, []byte{})
		hash = sha256.Sum256(blockInfo)
		temInt := big.Int{}
		temInt.SetBytes(hash[:])
		if temInt.Cmp(pow.target) == -1 {
			fmt.Printf("挖矿成功！hash:%x =====nonce:%d\n", hash[:], nonce)
			return hash[:], nonce
		} else {
			nonce++
		}
	}

}
