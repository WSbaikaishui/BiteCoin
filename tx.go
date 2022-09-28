/**
 * @Author: ZYW
 * @Description:  User Work
 * @File:  tx.go
 * @Date: 2022/9/28 15:10
 */
package main

import (
	"encoding/hex"
	"fmt"
	"log"
)

type TXOutput struct {
	Value int
	ScriptPubkey	string

}

type TXInput struct {
	Txid []byte
	Vout int
	ScriptSig	string
}

type Transaction struct {
	ID []byte
	Vin []TXInput
	Vout []TXOutput
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
					if out.CanBeUnlockedwith(address) && accumulated <amount{
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
	spentTXOs := 
}





func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

func (out *TXOutput) CanBelockedWith(unlockingData string) bool {
	return out.ScriptPubkey == unlockingData
}