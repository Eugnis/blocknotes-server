package root

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Web3NoteManager interface
type Web3NoteManager interface {
	InitClient() (*ethclient.Client, error)
	OpenWallet(privateKeyHex string) (*ecdsa.PrivateKey, *ecdsa.PublicKey)
	LoadBlock(client *ethclient.Client, blockNumber int64) (*types.Block, error)
}
