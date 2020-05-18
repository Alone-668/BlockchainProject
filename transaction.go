package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

const reward = 50.0
type Transaction struct {
	TXID []byte
	TXInputs []TXInput
	TXOutputs []TXOutput
}
type TXInput struct {
	TXid []byte
	Index int64
	//解锁脚本，我们用地址来模拟
	Sig string
}
type TXOutput struct {
	Value float64
	//锁定脚本,我们用地址模拟
	PubKeyHash string
}

//设置交易ID：交易ID为交易结构的hash
func (tx *Transaction)setHash()  {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(&tx)
	if err != nil {
		log.Panic("交易结构体编码失败")
	}
	data := buffer.Bytes()
	hash := sha256.Sum256(data)
	tx.TXID = hash[:]
}
func (tx *Transaction)IsCoinBase() bool {
	if len(tx.TXInputs)==0 && len(tx.TXInputs[0].TXid)==0 && tx.TXInputs[0].Index==-1 {
		return true
	}
	return false
}

func NewCoinBaseTX(address string,data string) *Transaction  {

	//挖矿交易的特点：
	//1. 只有一个input
	//2. 无需引用交易id
	//3. 无需引用index
	//矿工由于挖矿时无需指定签名，所以这个sig字段可以由矿工自由填写数据，一般是填写矿池的名字

	input := TXInput{TXid: []byte{}, Index: -1, Sig: data}
	output := TXOutput{reward, address}
	tx := Transaction{[]byte{}, []TXInput{input}, []TXOutput{output}}
	tx.setHash()
	return &tx
}
//创建普通的转账交易
//3. 创建outputs
//4. 如果有零钱，要找零

//func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	//map[2222] = []int64{0}
	//map[3333] = []int64{0, 1}
	utxOs, value := bc.FindNeedUTXOs(from, amount)
	if value < amount {
		fmt.Printf("余额不足，交易失败!")
		return nil
	}
	//2. 创建交易输入, 将这些UTXO逐一转成inputs
	var inputs []TXInput
	var outputs []TXOutput
	for TXID,UTXOArray:=range utxOs{
		for _,index :=range UTXOArray {
			input := TXInput{[]byte(TXID),int64(index),from}
			inputs = append(inputs, input)
		}
	}

	//创建交易输出
	output := TXOutput{amount, to}
	outputs = append(outputs,output)

	//找零
	if value>amount {
		outputs = append(outputs,TXOutput{value-amount,from})
	}
	tx := Transaction{[]byte{}, inputs, outputs}
	tx.setHash()
	return &tx
}