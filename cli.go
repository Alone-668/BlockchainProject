package main

import (
	"fmt"
	"os"
	"strconv"
)

type CLI struct {
	bc *BlockChain
}

const Usages = `
			printChain               "正向打印区块链"
			printChainR              "反向打印区块链"
			getBalance --address ADDRESS "获取指定地址的余额"
			send FROM TO AMOUNT MINER DATA "由FROM转AMOUNT给TO，由MINER挖矿，同时写入DATA"
			newWallet   "创建一个新的钱包(私钥公钥对)"
			listAddresses "列举所有的钱包地址"
`

func (cli *CLI)Run()  {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf(Usages)
		return
	}
	cmd := args[1]
	switch cmd {
	case "printChain":
		fmt.Printf("正向打印区块\n")
		cli.PrinBlockChain()
	case "printChainR":
		fmt.Printf("反向打印区块\n")
		cli.PrinBlockChainReverse()
	case "getBalance":
		fmt.Printf("获取余额\n")
		if len(args) == 4 && args[2] == "--address" {
			address := args[3]
			cli.GetBalance(address)
		}
	case "send":
		fmt.Printf("转账开始...\n")
		if len(args) != 7 {
			fmt.Printf("参数个数错误，请检查！\n")
			fmt.Printf(Usages)
			return
		}
		//./block send FROM TO AMOUNT MINER DATA "由FROM转AMOUNT给TO，由MINER挖矿，同时写入DATA"
		from := args[2]
		to := args[3]
		amount, _ := strconv.ParseFloat(args[4], 64) //知识点，请注意
		miner := args[5]
		data := args[6]
		cli.Send(from, to, amount, miner, data)
	case "newWallet":
		fmt.Printf("创建新的钱包...\n")
		cli.NewWallet()
	case "listAddresses":
		fmt.Printf("列举所有地址...\n")
		cli.ListAddresses()
	default:
		fmt.Printf("无效的命令，请检查!\n")
		fmt.Printf(Usages)
	}
	

}

