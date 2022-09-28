/**
 * @Author: ZYW
 * @Description:  User Work
 * @File:  blockchain
 * @Date: 2022/9/26 19:28
 */
package main

import (
	"fmt"
	"github.com/boltdb/bolt"
)

const dbFile = "blockchain_%s.db"
const blocksBucket = "blocks"

type BlockChain struct {
	tip []byte
	db *bolt.DB
}

type BlockchainIterator struct {
	currentHash []byte
	db *bolt.DB
}


func (bc *BlockChain) AddBlock(data string) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("1"))

		return nil
	}	)
	if err != nil{
		fmt.Println(err)
	}
	newBlock := NewBlock(data, lastHash)
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash,newBlock.Sertialize())
		err =  b.Put([]byte("1"),newBlock.Hash)
		if err != nil{
			fmt.Println(err)
		}
		bc.tip = newBlock.Hash


		return nil

	})

}

func NewGenesisBlock()	*Block  {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

func NewBlockchain() *BlockChain {
	var tip []byte

	db,err := bolt.Open(dbFile,0600,nil)
	if err != nil{
		fmt.Println(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil{
			gensis := NewGenesisBlock()
			b,err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil{
				fmt.Println(err)
			}
			err = b.Put(gensis.Hash,gensis.Sertialize())
			err = b.Put([]byte("1"),gensis.Hash)
			tip = gensis.Hash
		}else{
			tip = b.Get([]byte("1"))
		}
		return nil
	})
	bc := BlockChain{tip,db}
	return &bc

}

func (bc *BlockChain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip,bc.db}
	return bci

}

func (i *BlockchainIterator) Next() *Block {
	var block *Block
	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})
	if err != nil{
		fmt.Println(err)
	}
	i.currentHash = block.PrevBlockHash
	return block
}


//func main() {
//	bc := NewBlockchain()
//
//	bc.AddBlock("Send 1 BTC to Ivan")
//	bc.AddBlock("Send 2 more BTC to Ivan")
//
//	for _, block := range bc.blocks {
//		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
//		fmt.Printf("Data: %s\n", block.Data)
//		fmt.Printf("Hash: %x\n", block.Hash)
//		fmt.Println()
//	}
//
//	for _, block := range bc.blocks {
//		pow := NewPow(block)
//		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
//		fmt.Println()
//	}
//}

func main() {
	bc := NewBlockchain()
	defer bc.db.Close()
	cli := CLI{bc}
	cli.Run()
}