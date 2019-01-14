package notes

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

//Web3Manager ethereum token control
type Web3NoteManager struct{ Address string }

var deliminator = "||"

func (w *Web3NoteManager) InitClient() (*ethclient.Client, error) {
	client, err := ethclient.Dial(w.Address)
	return client, err
}

func (w *Web3NoteManager) OpenWallet(privateKeyHex string) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	return privateKey, publicKeyECDSA
}

func (w *Web3NoteManager) LoadBlock(client *ethclient.Client, blockNumber int64) (*types.Block, error) {
	block, err := client.BlockByNumber(context.Background(), big.NewInt(blockNumber))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	// log.Println(block.Hash())

	return block, nil
}

// func (w *Web3Manager) InitToken(tokenAddr string, client *ethclient.Client) (*SimlToken, error) {
// 	token, err := NewSimlToken(common.HexToAddress(tokenAddr), client)
// 	return token, err
// }

// func (w *Web3Manager) GetTokenData(tokenAddr string, client *ethclient.Client) error {
// 	token, err := NewSimlToken(common.HexToAddress(tokenAddr), client)
// 	if err != nil {
// 		log.Fatalf("Failed to instantiate a Token contract: %v", err)
// 	}
// 	name, err := token.Name(nil)
// 	symbol, err := token.Symbol(nil)
// 	totalSupply, err := token.TotalSupply(nil)

// 	if err != nil {
// 		log.Fatalf("Failed to retrieve token data: %v", err)
// 	}
// 	fmt.Println("Token name:", name)
// 	fmt.Println("Token symbol:", symbol)
// 	fmt.Println("Token totalSupply:", totalSupply)

// 	return err
// }
