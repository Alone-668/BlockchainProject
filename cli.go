package main

import (
	"fmt"
	"os"
)

type CLI struct {
	bc *BlockChain
}

const Usages = `
			addBlock --data DATA     "add data to blockchain"
			printChain               "print all blockchain data" 
`

func (cli *CLI)Run()  {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf(Usages)
		return
	}
	cmd := args[1]
	switch cmd {
	case "addBlock":
		fmt.Printf("添加区块\n")
		if len(args)==4 && args[2]=="--data"  {
			data := args[3]
			cli.AddBlock(data)
		}else {
			fmt.Printf("添加区块参数使用不当，请检查\n")
			fmt.Printf(Usages)
		}
	case "printChain":
		cli.PrinBlockChain()
	default:
		fmt.Printf("无效的命令，请检查!\n")
		fmt.Printf(Usages)
	}
	

}
