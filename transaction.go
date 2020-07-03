package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
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
//定义交易输入
type TXInput struct {
	//引用的交易ID
	TXid []byte
	//引用的output的索引值
	Index int64
	//解锁脚本，我们用地址来模拟
	//Sig string

	//真正的数字签名，由r，s拼成的[]byte
	Signature []byte

	//约定，这里的PubKey不存储原始的公钥，而是存储X和Y拼接的字符串，在校验端重新拆分（参考r,s传递）
	//注意，是公钥，不是哈希，也不是地址
	PubKey []byte
}

//定义交易输出
type TXOutput struct {
	//转账金额
	Value float64
	//锁定脚本,我们用地址模拟
	//PubKeyHash string

	//收款方的公钥的哈希，注意，是哈希而不是公钥，也不是地址
	PubKeyHash []byte
}
//由于现在存储的字段是地址的公钥哈希，所以无法直接创建TXOutput，
//为了能够得到公钥哈希，我们需要处理一下，写一个Lock函数
func (output *TXOutput) Lock(address string) {
	//1. 解码
	//2. 截取出公钥哈希：去除version（1字节），去除校验码（4字节）

	//真正的锁定动作！！！！！
	output.PubKeyHash = GetPubKeyFromAddress(address)
}

//给TXOutput提供一个创建的方法，否则无法调用Lock
func NewTXOutput(value float64, address string) *TXOutput {
	output := TXOutput{
		Value: value,
	}

	output.Lock(address)
	return &output
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
	//矿工由于挖矿时无需指定签名，所以这个PubKey字段可以由矿工自由填写数据，一般是填写矿池的名字
	//签名先填写为空，后面创建完整交易后，最后做一次签名即可
	input := TXInput{[]byte{}, -1, nil, []byte(data)}
	//output := TXOutput{reward, address}

	//新的创建方法
	output := NewTXOutput(reward, address)

	//对于挖矿交易来说，只有一个input和一output
	tx := Transaction{[]byte{}, []TXInput{input}, []TXOutput{*output}}
	tx.setHash()

	return &tx
}
//创建普通的转账交易
//3. 创建outputs
//4. 如果有零钱，要找零

//func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	//1. 创建交易之后要进行数字签名->所以需要私钥->打开钱包"NewWallets()"
	ws := NewWallets()

	//2. 找到自己的钱包，根据地址返回自己的wallet
	wallet := ws.WalletsMap[from]
	if wallet == nil {
		fmt.Printf("没有找到该地址的钱包，交易创建失败!\n")
		return nil
	}

	//3. 得到对应的公钥，私钥
	pubKey := wallet.PubKey
	privateKey := wallet.Private //稍后再用

	//传递公钥的哈希，而不是传递地址
	pubKeyHash := HashPubKey(pubKey)
	//map[2222] = []int64{0}
	//map[3333] = []int64{0, 1}
	utxOs, value := bc.FindNeedUTXOs(pubKeyHash, amount)
	if value < amount {
		fmt.Printf("余额不足，交易失败!")
		return nil
	}
	//2. 创建交易输入, 将这些UTXO逐一转成inputs
	var inputs []TXInput
	var outputs []TXOutput
	for TXID,UTXOArray:=range utxOs{
		for _,index :=range UTXOArray {
			input := TXInput{[]byte(TXID), int64(index), nil, pubKey}
			inputs = append(inputs, input)
		}
	}

	//创建交易输出
	output := NewTXOutput(amount, to)
	outputs = append(outputs, *output)

	//找零
	if value>amount {
		output = NewTXOutput(value-amount, from)
		outputs = append(outputs,*output)
	}
	tx := Transaction{[]byte{}, inputs, outputs}
	tx.setHash()


	bc.SignTransaction(tx,privateKey)


	return &tx
}
func (tx *Transaction)Sign(privateKey *ecdsa.PrivateKey, preTXs map[string]Transaction)  {

	txCopy:=tx.TrimmedCopy()
	//2. 循环遍历txCopy的inputs，得到这个input索引的output的公钥哈希,并赋值给要签名TX的公钥
	for i,input:=range txCopy.TXInputs{
		preTX := preTXs[string(input.TXid)]
		if len(preTX.TXID)==0 {
			log.Panic("引用的交易无效")
		}
		//先将输出TX复制（变成输入交易），然后把每个交易的signature和pubKey设为nil，
		//将输出交易的公钥hash赋值给 复制过得TX（输入交易） 的公钥字段
		txCopy.TXInputs[i].PubKey = preTX.TXOutputs[input.Index].PubKeyHash
		//3. 生成要签名的数据。要签名的数据一定是哈希值
		//a. 我们对每一个input都要签名一次，签名的数据是由当前input引用的output的哈希+当前的outputs（都承载在当前这个txCopy里面）
		//b. 要对这个拼好的txCopy进行哈希处理，SetHash得到TXID，这个TXID就是我们要签名最终数据。

		txCopy.setHash()
		//还原，以免影响后面input的签名
		txCopy.TXInputs[i].PubKey = nil
		signData := txCopy.TXID
		r, s, err := ecdsa.Sign(rand.Reader, privateKey, signData)
		if err != nil {
			log.Panic(err)
		}
		signNature:=append(r.Bytes(),s.Bytes()...)
		txCopy.TXInputs[i].Signature = signNature
	}

}

func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput
	for _,input :=range tx.TXInputs{
		inputs=append(inputs,TXInput{input.TXid,input.Index,nil,nil})
	}
	for _,output:=range tx.TXOutputs{
		outputs= append(outputs,output)
	}
	return Transaction{tx.TXID,inputs,outputs}
}