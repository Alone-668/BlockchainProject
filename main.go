package main

func main() {
	blockChain := NewBlockChain("张三")
	cli := CLI{blockChain}
	cli.Run()

}
