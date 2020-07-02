package main

func main() {
	blockChain := NewBlockChain("15Nfexuu4euziRE1gwzaBFVJ1WZcB7Dx5E")
	cli := CLI{blockChain}
	cli.Run()

}
