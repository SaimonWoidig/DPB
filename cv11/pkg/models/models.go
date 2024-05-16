// Code generated by "gocqlx/cmd/schemagen"; DO NOT EDIT.

package models

import (
	"github.com/scylladb/gocqlx/v2/table"
	"time"
)

// Table models.
var (
	Messages = table.New(table.Metadata{
		Name: "messages",
		Columns: []string{
			"message",
			"room_id",
			"speaker_id",
			"time",
		},
		PartKey: []string{
			"room_id",
		},
		SortKey: []string{
			"time",
		},
	})
)

type MessagesStruct struct {
	Message   string
	RoomId    int64
	SpeakerId int64
	Time      time.Time
}