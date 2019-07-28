package handle

import (
	"encoding/json"
	"github.com/davyxu/cellnet"
	"socket_server/constant"
	"socket_server/server"
	"socket_server/server/message"
	"time"
)

var UnReadMessageHandle = &server.Handler{Run: func(genericPeer cellnet.GenericPeer, event cellnet.Event) error {
	// 获取未读消息内容
	getUnReadMsg := event.Message().(*message.Message)
	unReadMsgList, err := server.GetMsgSync(getUnReadMsg)
	if err != nil {
		return err
	}
	unReadMsgListString, err := json.Marshal(unReadMsgList)
	if err != nil {
		return err
	}
	// 将未读消息返回给客户端
	event.Session().Send(&message.Message{
		From:       constant.ServerId,
		To:         getUnReadMsg.From,
		Cmd:        message.CommandType_GET_UN_READ_MESSAGE_RESP,
		CreateTime: time.Now().Unix() / 1e6,
		MsgType:    message.MsgType_TEXT,
		ChatType:   message.ChatType_PRIVATE,
		GroupId:    constant.Empty,
		Content:    string(unReadMsgListString),
		Extras:     constant.Empty,
	})
	return nil
}}
