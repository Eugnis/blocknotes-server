package mongo

import (
	root "blocknotes_server/pkg"
	"blocknotes_server/pkg/notes"
	"context"
	"fmt"
	"log"
	"time"
	"unicode/utf8"

	"github.com/ethereum/go-ethereum/common"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/h2non/filetype"
)

type NoteService struct {
	collection *mgo.Collection
}

func NewNoteService(session *Session, dbName string, collectionName string) *NoteService {
	collection := session.GetCollection(dbName, collectionName)
	collection.EnsureIndex(noteModelIndex())
	return &NoteService{collection}
}

func (p *NoteService) Create(u *root.Note) error {
	note := newNoteModel(u)
	return p.collection.Insert(&note)
}

func (p *NoteService) GetByNoteID(id string) (*root.Note, error) {
	model := noteModel{}
	err := p.collection.FindId(bson.ObjectIdHex(id)).One(&model)
	return model.toRootNote(), err
}

func (p *NoteService) GetByNoteAddress(address string) (*root.Note, error) {
	model := noteModel{}
	err := p.collection.Find(bson.M{"hash": address}).One(&model)
	return model.toRootNote(), err
}

func (p *NoteService) ListNotes(nsr root.NoteSearch) ([]*root.Note, error) {
	models := []noteModel{}
	searchReq := bson.M{}
	gibberish := []string{"gVxX4N7rd0", "coinbenerefuel", "Ignore", "hotwallet drain fee", "BFX_REFILL_SWEEP", "service charge"}
	if nsr.NetName != "" {
		searchReq["netname"] = nsr.NetName
		// log.Println(nsr.NetName)
	}
	if nsr.NetType != "" {
		searchReq["nettype"] = nsr.NetType
		// log.Println(nsr.NetType)
	}
	if nsr.DataType != "" {
		searchReq["datatype"] = bson.M{"$regex": nsr.DataType, "$options": "i"}
		// log.Println(nsr.DataType)
	} else {
		searchReq["datatype"] = bson.M{"$ne": "text"}
	}
	if nsr.SearchText != "" && nsr.SearchType != "" {
		searchReq[nsr.SearchType] = bson.M{"$regex": nsr.SearchText, "$options": "i"}
		// log.Println(searchReq[nsr.SearchType])
	} else {
		searchReq["textpreview"] = bson.M{"$nin": gibberish}
	}

	q := p.collection.Find(searchReq).Sort("-$natural")
	q = q.Skip(nsr.From).Limit(nsr.Count)

	err := q.All(&models)
	if err != nil {
		fmt.Println("error ", err)
	}
	result := make([]*root.Note, len(models))
	for i, item := range models {
		result[i] = item.toRootNote()
		// log.Println(item.DataSize)
	}
	return result, err
}

func (p *NoteService) Update(u *root.Note) error {
	model, err := p.GetByNoteID(u.ID.Hex())
	if err != nil {
		return err
	}

	if u.NetName != "" {
		model.NetName = u.NetName
	}
	if u.NetType != "" {
		model.NetType = u.NetType
	}
	if u.Address != "" {
		model.Address = u.Address
	}
	if u.BlockNum != model.BlockNum {
		model.BlockNum = u.BlockNum
	}
	if u.DataType != "" {
		model.DataType = u.DataType
	}
	if u.DataSize != model.DataSize {
		model.DataSize = u.DataSize
	}
	if u.TextPreview != "" {
		model.TextPreview = u.TextPreview
	}
	if u.TxTime != model.TxTime {
		model.TxTime = u.TxTime
	}
	// log.Printf(u.MainAttributeID, model.MainAttributeID)

	return p.collection.UpdateId(bson.ObjectIdHex(model.ID.Hex()), model)
}

func (p *NoteService) Remove(id string) error {
	err := p.collection.RemoveId(bson.ObjectIdHex(id))
	return err
}

func (p *NoteService) NotesFetcher() {
	ticker := time.NewTicker(60 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("Current task minute")

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (p *NoteService) BatchNotesFetcher() {
	for {
		const nodeAddr = "https://mainnet.infura.io/v3/72972d2f7624494497fdf19174fd1083"
		// const localNodeAddr = "ws://127.0.0.1:8546"

		var lastRecord root.Note
		var startBlock int64
		err := p.collection.Find(bson.M{}).Sort("-$natural").Limit(1).One(&lastRecord)
		if err != nil {
			log.Println("Error fetching notes", err.Error())
			startBlock = 0
		} else {
			startBlock = lastRecord.BlockNum + 1
		}
		// startBlock = 122760
		infuraManager := notes.Web3NoteManager{Address: nodeAddr}
		// localManager := notes.Web3NoteManager{Address: localNodeAddr}
		infuraClient, infuraErr := infuraManager.InitClient()
		if infuraErr != nil {
			log.Println("Error connecting infura node", infuraErr.Error())
			return
		}
		// localClient, localErr := infuraManager.InitClient()
		// if localErr != nil {
		// 	log.Println("Error connecting local node", localErr.Error())
		// 	return
		// }

		header, err := infuraClient.HeaderByNumber(context.Background(), nil)
		if err != nil {
			log.Fatal(err)
		}
		lastBlock := header.Number.Int64()

		log.Printf("Start: %d, end: %d", startBlock, lastBlock)

		for curBlock := startBlock; curBlock <= lastBlock; curBlock++ {
			bulk := p.collection.Bulk()
			foundCount := 0
			block, errBlock := infuraManager.LoadBlock(infuraClient, curBlock)
			if errBlock != nil {
				log.Println("Error block", errBlock.Error())
			}
			for _, tx := range block.Transactions() {
				if len(tx.Data()) > 3 {
					var newNote root.Note
					notePresent := false
					isContract := false
					addrTo := ""
					// if len(hexEnc) > 10 && hexEnc[10:11] == "0" {
					// 	// log.Println("Contract call!", tx.Hash().String())
					// 	isContract = true
					// }
					if utf8.Valid(tx.Data()) {
						encodedString := string(tx.Data())
						if len(encodedString) > 1 && encodedString != "fz/X" && encodedString != "undefined" && tx.Data()[0] != 0 {
							if tx.To() != nil {
								btcode, btcodeerr := infuraClient.CodeAt(context.Background(), *tx.To(), nil)
								if btcodeerr != nil {
									log.Println("err", btcodeerr.Error())
								} else {
									if len(btcode) > 0 {
										isContract = true
									}
								}
								addrTo = tx.To().String()
							}
							// fmt.Printf("Hash: %s Len: %d String: '%s'\n", tx.Hash().String(), len(encodedString), encodedString)
							if !isContract {
								newNote = root.Note{
									ID:          bson.NewObjectId(),
									NetName:     "ethereum",
									NetType:     "mainnet",
									Hash:        tx.Hash().String(),
									Address:     addrTo,
									BlockNum:    curBlock,
									DataType:    "text",
									DataSize:    len(tx.Data()),
									TextPreview: encodedString,
									TxTime:      time.Unix(block.Time().Int64(), 0)}
								notePresent = true
							}

						}
					} else {
						kind, unknown := filetype.Match(tx.Data())
						if unknown == nil && kind.Extension != "unknown" {
							// fmt.Printf("Hash: %s Len: %d File type: %s. MIME: %s\n", tx.Hash().String(), len(tx.Data()), kind.Extension, kind.MIME.Value)
							addrTo := ""
							if tx.To() != nil {
								addrTo = tx.To().String()
							}
							newNote = root.Note{
								ID:          bson.NewObjectId(),
								NetName:     "ethereum",
								NetType:     "mainnet",
								Hash:        tx.Hash().String(),
								Address:     addrTo,
								BlockNum:    curBlock,
								DataType:    kind.MIME.Value,
								DataSize:    len(tx.Data()),
								TextPreview: kind.Extension,
								TxTime:      time.Unix(block.Time().Int64(), 0)}
							notePresent = true
						}
					}

					if notePresent {
						foundCount++
						bulk.Insert(&newNote)
					}

				}
			}

			if foundCount > 0 {
				log.Printf("[%d] Notes: %d (%s)", curBlock, foundCount, time.Unix(block.Time().Int64(), 0).UTC())
				_, runned := bulk.Run()
				if runned != nil {
					log.Println("Error bulk", runned.Error())
					// return
				}
			}
		}
		log.Println("Block scan finished, restarting in 5 minutes")
		time.Sleep(5 * time.Minute)
	}
	// p.BatchNotesFetcher()
}

func (p *NoteService) NotesFixer() {
	const nodeAddr = "https://mainnet.infura.io/v3/72972d2f7624494497fdf19174fd1083"
	infuraManager := notes.Web3NoteManager{Address: nodeAddr}
	// localManager := notes.Web3NoteManager{Address: localNodeAddr}
	infuraClient, infuraErr := infuraManager.InitClient()
	if infuraErr != nil {
		log.Println("Error connecting infura node", infuraErr.Error())
		return
	}

	models := []noteModel{}
	bulk := p.collection.Bulk()
	err := p.collection.Find(bson.M{"textpreview": ""}).All(&models)
	if err != nil {
		log.Fatal(err)
		return
	}
	cnt := 0
	log.Printf("Fixing %d models", len(models))
	for _, note := range models {
		tx, _, err := infuraClient.TransactionByHash(context.Background(), common.HexToHash(note.Hash))
		if err != nil {
			log.Fatal(err)
			return
		}
		// if len(hexEnc) > 10 && hexEnc[10:11] == "0" {
		// 	// log.Println("Contract call!", tx.Hash().String())
		// 	isContract = true
		// }
		kind, unknown := filetype.Match(tx.Data())
		if unknown == nil && kind.Extension != "unknown" {
			// fmt.Printf("Hash: %s Len: %d File type: %s. MIME: %s\n", tx.Hash().String(), len(tx.Data()), kind.Extension, kind.MIME.Value)

			bulk.Update(bson.M{"_id": note.ID}, bson.M{"$set": bson.M{"datatype": kind.MIME.Value, "textpreview": kind.Extension}})
		}
	}

	res, runned := bulk.Run()
	if runned != nil {
		log.Println("Error bulk", runned.Error())
		return
	}
	log.Printf("Matched %d, Modified %d, Fixed %d", res.Matched, res.Modified, cnt)

}
