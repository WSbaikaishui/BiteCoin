/**
 * @Author: ZYW
 * @Description:  User Work
 * @File:  block
 * @Date: 2022/9/26 16:37
 */
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

type Block struct {
	Timestamp int64
	Transactions     []*Transaction
	PrevBlockHash	[]byte
	Hash     []byte
	Nonce int
}
//type BlockHeader struct {
//	version int32
//	PrevBlock  chainhash.hash
//	MerkleRoot chainhash.hash
//	TimeStamp time.Time
//	Bits uint32
//	Nonce	uint32
//
//}

//func (b *Block)SetHash()  {
//	timestamp:= []byte(strconv.FormatInt(b.Timestamp,10))
//	headers := bytes.Join([][]byte{b.PrevBlockHash,b.Data,timestamp},[]byte{})
//	hash:= sha256.Sum256(headers)
//	//TODO 为啥要这么赋值
//	//直接b.hash = hash不行吗
//	b.Hash = hash[:]
//}

func NewBlock(tx []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(),tx,prevBlockHash,[]byte{},0}
	pow := NewPow(block)
	nonce ,hash := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return block

}

func (b *Block) Sertialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err !=nil {
		fmt.Println(err)
		}
		return result.Bytes()
}

func DeserializeBlock(d []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil{
		fmt.Println(err)
	}
	return &block
}








