package main

import "fmt"

func main() {
	blockChain := NewBlockChain()
	blockChain.AddBlock("李琪给龚睿50比特币")
	blockChain.AddBlock("李琪又给龚睿50比特币")
	for i, block := range blockChain.blocks {
		fmt.Printf("======== 当前区块高度： %d ========\n", i)
		fmt.Printf("前区块哈希值： %x\n", block.PrevHash)
		fmt.Printf("当前区块哈希值： %x\n", block.Hash)
		fmt.Printf("区块数据 :%s\n", block.Data)
	}

}
