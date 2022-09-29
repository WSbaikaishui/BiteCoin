/**
 * @Author: ZYW
 * @Description:  User Work
 * @File:  tx.go
 * @Date: 2022/9/28 15:10
 */
package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
)

type TXOutput struct {
	Value int
	PubKeyHash	[]byte

}

type TXInput struct {
	Txid []byte
	Vout int
	Signature []byte
	PubKey []byte
}

type Transaction struct {
	ID []byte
	Vin []TXInput
	Vout []TXOutput
}

// UsesKey 方法检查输入使用了指定密钥来解锁一个输出。
//输入存储的是原生的公钥（也就是没有被哈希的公钥）
//但是这个函数要求的是哈希后的公钥。
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)
	return bytes.Compare(lockingHash,pubKeyHash) == 0
}

func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash) - 4]
	out.PubKeyHash = pubKeyHash
}

func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash,pubKeyHash) == 0
}


func NewCoinbaseTX(to, data string) *Transaction {
	if data == ""{
		data = fmt.Sprintf("Reward to '%s'",to)
	}
	txin := TXInput{[]byte{},-1,data}
	txout := TXOutput{subsidy,to}
	tx := Transaction{nil,[]TXInput{txin},[]TXOutput{txout}}
	tx.setID()
	return &tx
}

func NewUTXOTransaction(from, to string, amount int, bc *BlockChain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput


	acc, validOutputs := bc.FindSpendableOutputs(from,amount)
	if acc <amount{
		log.Panic("error: not enough funds")
		}
		for txid, outs :=range validOutputs{
			txID,err := hex.DecodeString(txid)
			for _,out := range outs{
				input := TXInput{txID,out,from}
				inputs = append(inputs,input)
			}
		}

		outputs = append(outputs,TXOutput{amount,to})
		if acc > amount{
			outputs = append(outputs,TXOutput{acc -amount,from})
		}
		tx := Transaction{nil,inputs,outputs}
		tx.SetID()

		return &tx
}

func (bc *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(address)
	accumulated := 0

	Work :
			for _,tx := range unspentTXs{
				txID := hex.EncodeToString(tx.ID)

				for outIdx,out := range tx.Vout{
					if out.CanBeUnlockedWith(address) && accumulated <amount{
						accumulated +=out.Value
						unspentOutputs[txID] = append(unspentOutputs[txID],outIdx)

						if accumulated >= amount{
							break Work
						}
					}
				}

			}
			return accumulated, unspentOutputs
}

func (bc *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspentTXs []Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()
	for {
		block := bci.Next()
		for _,tx := range block.Transactions{
			txID := hex.EncodeToString(tx.ID)

			Outputs:
				for outIdx,out := range tx.Vout{
					if spentTXOs[txID] !=nil{
						for _,spentOut := range spentTXOs[txID]{
							if spentOut == outIdx{
								continue Outputs
							}
						}
					}
					if out.CanBeUnlockedWith(address){
						unspentTXs = append(unspentTXs,*tx)
					}
				}

				if tx.IsCoinbase() == false{
					for _,in := range tx.Vin{
						if in.CanUnlockOutputWith(address){
							inTxID := hex.EncodeToString(in.Txid)
							spentTXOs[inTxID] = append(spentTXOs[inTxID],in.Vout)
						}
					}
				}
		}

		if len(block.PrevBlockHash) == 0{
			break
		}
	}
	return unspentTXs
}





func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubkey == unlockingData
}