/*
 * This is a 'Toy Blockchain' created to learn basic concepts
 * behind a Blockchain. It runs on a single node, and is in memory.
 * For understanding the terminologies refer:
 * https://ethereum.org/en/developers/docs/intro-to-ethereum/#terminology
 */
package main

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"
)

// Max number of transactions to be packed in a Block
const MAX_TXNS_PER_BLOCK = 5

type Transaction struct {
	payer string
	payee string
	amt   float64
}

type Block struct {
	data     []Transaction // list of transactions in the Block
	prevHash string        // hash of the previous Block
	unixTs   int64         // unix timestamp when the Block was created
	nonce    int           // Proof Of Work
	hash     string        // hash of the Block
}

type BlockChain struct {
	current    *Block  // Current Block for outstanding transactions
	chain      []Block // Committed Blocks
	difficulty int     // Proof Of Work difficulty
}

// Cryptographic Hash using SHA-256
func SHA256(packedBytes []byte) string {
	hash := sha256.Sum256(packedBytes)
	return fmt.Sprintf("%x", hash)
}

// Proof Of Work
func (b *Block) mine(difficulty int) {
	fixedBlockBytes := []byte(fmt.Sprintf("%v", b.data) + fmt.Sprintf("%v", b.prevHash) + fmt.Sprintf("%v", b.unixTs))
	for !strings.HasPrefix(b.hash, strings.Repeat("0", difficulty)) {
		b.nonce++
		b.hash = SHA256(append(fixedBlockBytes, []byte(fmt.Sprintf("%v", b.nonce))...))
	}
}

func (b Block) PrettyDisplay() {
	fmt.Print("\n\nBlock: ")
	for _, txn := range b.data {
		fmt.Printf("\n%+v", txn)
	}
	fmt.Printf("\nnonce: %v", b.nonce)
	fmt.Printf("\nprevHash: %v", b.prevHash)
	fmt.Printf("\nunixTimestamp: %v", b.unixTs)
	fmt.Printf("\nHash: %v", b.hash)
	fmt.Print("\n\t\t|\n\t\t|\n\t\tv")
}

func CreateBlockChain(difficulty int) BlockChain {
	genesisBlock := Block{
		unixTs: time.Now().UnixMicro(),
		nonce:  0,
	}
	genesisBlock.mine(difficulty)
	bc := BlockChain{
		current:    nil,
		chain:      []Block{genesisBlock},
		difficulty: difficulty,
	}
	return bc
}

func (bc BlockChain) lastBlock() *Block {
	return &bc.chain[len(bc.chain)-1]
}

func (bc *BlockChain) AddTxn(txn Transaction) {
	if bc.current == nil || len(bc.current.data) >= MAX_TXNS_PER_BLOCK {
		bc.CommitBlock()
		bc.newBlock(txn)
	} else {
		// Append txn to current block in the BlockChain
		bc.current.data = append(bc.current.data, txn)
	}
}

/*
 * Create a new Block and add a transaction
 * Calling this without previously calling CommitBlock will discard the
 * outstanding transactions
 */
func (bc *BlockChain) newBlock(txn Transaction) {
	bc.current = &Block{
		data:     []Transaction{txn},
		prevHash: bc.lastBlock().hash,
		unixTs:   time.Now().UnixMicro(),
	}
}

/*
 * Commit outstanding transactions in the current Block,
 * and append the current Block to the BlockChain
 */
func (bc *BlockChain) CommitBlock() {
	if bc.current != nil {
		bc.current.mine(bc.difficulty)
		bc.chain = append(bc.chain, *bc.current)
		bc.current = nil
	}
}

func (bc BlockChain) PrettyDisplay() {
	fmt.Println("\n--------- BlockChain Start -----------")
	fmt.Printf("Proof Of Work Diffculty: %v (no. of leading 0s in the hash)", bc.difficulty)
	for _, b := range bc.chain {
		b.PrettyDisplay()
	}
	fmt.Print("\n\n--------- BlockChain End -----------\n\n")
}

func main() {
	blockchain := CreateBlockChain(4)

	// Simulate adding transactions
	blockchain.AddTxn(Transaction{
		payer: "alice",
		payee: "bob",
		amt:   10.0,
	})
	blockchain.AddTxn(Transaction{
		payer: "alice",
		payee: "bob",
		amt:   30.0,
	})
	blockchain.AddTxn(Transaction{
		payer: "bob",
		payee: "alice",
		amt:   35.0,
	})
	blockchain.AddTxn(Transaction{
		payer: "clark",
		payee: "bob",
		amt:   10.0,
	})
	blockchain.AddTxn(Transaction{
		payer: "clark",
		payee: "alice",
		amt:   5.0,
	})
	blockchain.AddTxn(Transaction{
		payer: "clark",
		payee: "bob",
		amt:   10.0,
	})
	blockchain.AddTxn(Transaction{
		payer: "clark",
		payee: "alice",
		amt:   5.0,
	})

	// Commit outstanding transactions if the last block is not full
	blockchain.CommitBlock()
	blockchain.PrettyDisplay()
}
