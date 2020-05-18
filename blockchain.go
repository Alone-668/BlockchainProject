package main

import (
	"BlockchainProject/bolt"
	"bytes"
	"fmt"
	"log"
	"time"
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
func NewBlockChain(address string) *BlockChain {
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
			genesisBlock := GenesisBlock(address)
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
func GenesisBlock(address string) *Block {
	coinBase := NewCoinBaseTX(address, "创世块")
	return NewBlock([]*Transaction{coinBase}, []byte{})
}

//5. 添加区块
func (bc *BlockChain) AddBlock(txs []*Transaction) {

	db := bc.db
	lastHash := bc.tail

	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Blockbucket))
		if bucket == nil {
			log.Panic("blockBucket不应该为空，bolt数据库出现问题")
		}
		block := NewBlock(txs, lastHash)
		bucket.Put(block.Hash,block.Serialize())
		bucket.Put([]byte("LastHashKey"),block.Hash)
		bc.tail = block.Hash

		return nil
	})

}

func (bc *BlockChain) Printchain() {

	blockHeight := 0
	bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Blockbucket))
		//从第一个key-> value 进行遍历，到最后一个固定的key时直接返回
		bucket.ForEach(func(k, v []byte) error {
			if bytes.Equal(k,[]byte("LastHashKey")) {
				return nil
			}
			block := Deserialize(v)
			fmt.Printf("=============== 区块高度: %d ==============\n", blockHeight)
			blockHeight++
			fmt.Printf("版本号: %d\n", block.Version)
			fmt.Printf("前区块哈希值: %x\n", block.PrevHash)
			fmt.Printf("梅克尔根: %x\n", block.MerkelRoot)
			timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
			fmt.Printf("时间戳: %s\n", timeFormat)
			fmt.Printf("难度值(随便写的）: %d\n", block.Difficulty)
			fmt.Printf("随机数 : %d\n", block.Nonce)
			fmt.Printf("当前区块哈希值: %x\n", block.Hash)
			//只有coinbase有数据
			fmt.Printf("区块数据 :%s\n", block.Transactions[0].TXInputs[0].Sig)
			return nil

		})
		return nil
	})
}
//找到指定地址的所有的utxo
func (bc *BlockChain)FindUTXOs(address string) []TXOutput  {
	var UTXO []TXOutput
	txs := bc.FindUTXOTransactions(address)
	for _,tx:=range txs{
		for _,output :=range tx.TXOutputs{
			if output.PubKeyHash == address {
				UTXO = append(UTXO, output)
			}
		}
	}
	return UTXO
}
//根据需求找到合理的utxo
func (bc *BlockChain)FindNeedUTXOs(from string,value float64) (map[string][]int64,float64) {
	utxos:=make(map[string][]int64)
	var calc float64
	txs := bc.FindUTXOTransactions(from)
	for _,tx:=range txs{
		for i,output:=range tx.TXOutputs{
			if output.PubKeyHash == from {
				if calc<value {
					//不满足交易额度
					//2. 统计一下当前utxo的总额
					//第一次进来: calc =3,  map[3333] = []uint64{0}
					//第二次进来: calc =3 + 2,  map[3333] = []uint64{0, 1}
					//第三次进来：calc = 3 + 2 + 10， map[222] = []uint64{0}
					utxos[string(tx.TXID)] = append(utxos[string(tx.TXID)],int64(i))
					calc+=output.Value
					if calc>value {
						fmt.Printf("找到了满足的金额：%f\n", calc)
						return utxos,calc
					}else {
						fmt.Printf("%s 不满足转账金额,当前总额：%f， 目标金额: %f\n",from, calc, value)
					}
				}
			}
		}
	}
	return utxos,calc
}


func (bc *BlockChain)FindUTXOTransactions(address string) []*Transaction {
	var txs []*Transaction
	//定义一个消耗过的utxomap集合
	SpendUTXOMap:=make(map[string][]int64)
	it := bc.NewItertor()
	for  {
		block := it.Next()

		for _,tx:=range block.Transactions {

			OUTPUT:
			for i,ouput:=range tx.TXOutputs{
				//map[2222] = []int64{0}
				//map[3333] = []int64{0, 1}
				if SpendUTXOMap[string(tx.TXID)] != nil {
					for  _,j :=range SpendUTXOMap[string(tx.TXID)]{
						if  int64(i) == j{
							continue OUTPUT
						}
					}
				}

				if ouput.PubKeyHash == address {
					txs = append(txs,tx)
				}

			}

			//遍历inputs
			if !tx.IsCoinBase() {
				for _,input:=range tx.TXInputs{
					//判断是否是发起交易的人的sig签名
					if input.Sig == address{
						SpendUTXOMap[string(input.TXid)] = append(SpendUTXOMap[string(input.TXid)],input.Index)
					}
				}
			}

		}
		if len(block.PrevHash) == 0 {
			break
			fmt.Printf("区块链遍历结束")
		}
	}
	return txs
}
