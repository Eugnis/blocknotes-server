package mongo

import (
	root "blocknotes_server/pkg"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type noteModel struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	NetName     string        `json:"net_name,omitempty"`
	NetType     string        `json:"net_type,omitempty"`
	Hash        string        `json:"hash,omitempty"`
	Address     string        `json:"address,omitempty"`
	BlockNum    int64         `json:"block_num,omitempty"`
	DataType    string        `json:"data_type,omitempty"`
	DataSize    int           `json:"data_size,omitempty"`
	TextPreview string        `json:"text_preview,omitempty"`
	TxTime      time.Time     `json:"tx_time,omitempty"`
}

func noteModelIndex() mgo.Index {
	return mgo.Index{
		Key:        []string{"hash"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
}

func newNoteModel(u *root.Note) *noteModel {
	return &noteModel{
		ID:          bson.ObjectId(u.ID),
		NetName:     u.NetName,
		NetType:     u.NetType,
		Hash:        u.Hash,
		Address:     u.Address,
		BlockNum:    u.BlockNum,
		DataType:    u.DataType,
		DataSize:    u.DataSize,
		TextPreview: u.TextPreview,
		TxTime:      u.TxTime}
}

func (u *noteModel) toRootNote() *root.Note {
	return &root.Note{
		ID:          u.ID,
		NetName:     u.NetName,
		NetType:     u.NetType,
		Hash:        u.Hash,
		Address:     u.Address,
		BlockNum:    u.BlockNum,
		DataType:    u.DataType,
		DataSize:    u.DataSize,
		TextPreview: u.TextPreview,
		TxTime:      u.TxTime}
}
