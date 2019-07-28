package models

import (
	"socket_server/web/database"
	"time"
)

type Message struct {
	MessageID  uint64    `xorm:"pk autoincr notnull 'message_id'"`
	FromId     uint64    `xorm:"bigint(20) 'from_id'"`
	ToId       uint64    `xorm:"bigint(20) 'to_id'"`
	CreateTime time.Time `xorm:"create_time"`
	MsgType    uint32    `xorm:"int(11) 'msg_type'"`
	ChatType   uint32    `xorm:"int(11) 'chat_type'"`
	GroupId    string    `xorm:"varchar(200) 'group_id'"`
	Content    string    `xorm:"text 'content'"`
	Extras     string    `xorm:"varchar(200) 'extras'"`
}

func SaveMessage(message *Message) error {
	if _, err := database.DB.
		Insert(message); err != nil {
		return err
	}
	return nil
}
