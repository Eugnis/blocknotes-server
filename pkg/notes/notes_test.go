package notes_test

import (
	"blocknotes_server/pkg/notes"
	"context"
	"fmt"
	"log"
	"testing"
	"unicode/utf8"

	"github.com/h2non/filetype"
)

// const nodeAddr = "https://mainnet.infura.io/v3/72f3efb50f7942009d064766215cb2d5"

const nodeAddr = "ws://127.0.0.1:8546"

const privateKey = "443C45DB6F9F6EB1E4AA3D995E9D245DCF9FF40CF6460AEA50D468FA694F9498"
const tokenAddr = "0x2491b76e89c3da92c4e13b12a0ba4fee37c25e53"

func Test_Notes(t *testing.T) {
	t.Run("Notes read", web3_client_test)
	// t.Run("Balance of address", token_current_balance)
}

func data_test(t *testing.T) {
	if len(tx.Data()) > 3 {
		// fmt.Println("txfat", tx.Data()[0])
		if utf8.Valid(tx.Data()) {
			encodedString := string(tx.Data())
			// fmt.Printf("Hash: %s Len: %d String: '%s'\n", tx.Hash().String(), len(encodedString), encodedString)
			if len(encodedString) > 1 && tx.Data()[0] != 0 {
				fmt.Printf("111Hash: %s Len: %d Byte: %d String: '%s'\n", tx.Hash().String(), len(encodedString), tx.Data()[0], encodedString)
			}
		} else {
			kind, unknown := filetype.Match(tx.Data())
			if unknown == nil && kind.Extension != "unknown" {
				fmt.Printf("Hash: %s Len: %d File type: %s. MIME: %s\n", tx.Hash().String(), len(tx.Data()), kind.Extension, kind.MIME.Value)
			}
		}

	}
}

func web3_client_test(t *testing.T) {
	//Arrange
	c := notes.Web3NoteManager{}
	client, initClientErr := c.InitClient(nodeAddr)

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	last_block := header.Number.Int64()
	// last_block := int64(6985440)
	from_block := last_block - 100

	fmt.Printf("From %d To %d\n", from_block, last_block) // 5671744

	for j := from_block; j <= last_block; j++ {
		fmt.Printf("Loading block %d\n", j)
		block, _ := c.LoadBlock(client, j)
		for _, tx := range block.Transactions() {
			if len(tx.Data()) > 3 {
				// fmt.Println("txfat", tx.Data()[0])
				if utf8.Valid(tx.Data()) {
					encodedString := string(tx.Data())
					// fmt.Printf("Hash: %s Len: %d String: '%s'\n", tx.Hash().String(), len(encodedString), encodedString)
					if len(encodedString) > 1 && tx.Data()[0] != 0 {
						fmt.Printf("111Hash: %s Len: %d Byte: %d String: '%s'\n", tx.Hash().String(), len(encodedString), tx.Data()[0], encodedString)
					}
				} else {
					kind, unknown := filetype.Match(tx.Data())
					if unknown == nil && kind.Extension != "unknown" {
						fmt.Printf("Hash: %s Len: %d File type: %s. MIME: %s\n", tx.Hash().String(), len(tx.Data()), kind.Extension, kind.MIME.Value)
					}
				}

			}
		}
	}

	//Act
	//private, public := c.OpenWallet(privateKey)
	// token, initTokenErr := c.InitToken(tokenAddr, client)
	// name, callingError := token.Name(nil)

	// token.Nam(nil)
	t.Error("Done")
	//Assert
	if initClientErr != nil {
		t.Error("Error InitClient")
	}
	// if loadBlockErr != nil {
	// 	t.Error("loadBlockErr error")
	// }
	// if public == nil {
	// 	t.Error("Error open wallet")
	// }
	// if callingError != nil {
	// 	t.Error("Token function calling error")
	// }
	// if tokenName != name {
	// 	t.Error("Error token data receiving")
	// }
}
