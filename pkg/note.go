package root

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// Note struct
type Note struct {
	ID          bson.ObjectId `json:"id,omitempty" bson:"_id"`
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

// NoteService interface
type NoteService interface {
	Create(u *Note) error
	GetByNoteID(id string) (*Note, error)
	GetByNoteAddress(address string) (*Note, error)
	ListNotes(nsr NoteSearch) ([]*Note, int, error)
	Update(s *Note) error
	Remove(id string) error
	NotesFetcher()
	BatchNotesFetcher()
	NotesFixer()
}
