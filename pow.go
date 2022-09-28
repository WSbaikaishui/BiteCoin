/**
 * @Author: ZYW
 * @Description:  User Work
 * @File:  pow
 * @Date: 2022/9/26 20:15
 */
package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strings"
)

const targetBits = 24
const maxNonce =  math.MaxInt64
type Pow struct {
	block *Block
	target *big.Int
}

func NewPow(b *Block) *Pow {
	target := big.NewInt(1)
	target.Lsh(target,uint(256-targetBits))
	pow := &Pow{b,target}
	return pow
}

func (pow *Pow) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.HashTransactions(),
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
		)
	return data

}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash	[32]byte
	for _,tx := range b.Transactions{
		txHashes = append(txHashes,tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes,[]byte{}))
	return txHash[:]
}


func (pow *Pow)Run() (int, []byte)  {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining this block containing \"%s\"\n", pow.block.Data)

	for nonce <maxNonce{
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1{
			fmt.Printf("\r%x",hash)
			break
		}else {
			nonce++
		}

	}
	fmt.Print("\n\n")
	return nonce , hash[:]


}

func (pow *Pow) Validate() bool {
	var hashInt big.Int
	data := pow.prepareData(pow.block.Nonce)
	hash :=sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	return hashInt.Cmp(pow.target) == -1
}




func IntToHex(ten int64) []byte {
	m := int64(0)
	hex := make([]int64, 0)
	for {
		m = ten % 16
		ten = ten / 16
		if ten == 0 {
			hex = append(hex, m)
			break
		}
		hex = append(hex, m)
	}
	hexStr := []string{}
	for i:=len(hex)-1;i>=0;i--{
		if hex[i] >= 10 {
			hexStr = append(hexStr, fmt.Sprintf("%c", 'A'+hex[i]-10))
		} else {
			hexStr = append(hexStr, fmt.Sprintf("%d", hex[i]))
		}
	}
	return []byte(strings.Join(hexStr, ""))
}
