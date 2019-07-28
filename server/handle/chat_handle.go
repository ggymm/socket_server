package handle

import (
	"github.com/davyxu/cellnet"
	"socket_server/server"
	"socket_server/server/message"
	"socket_server/utils"
	"strconv"
	"time"
)

var ChatHandle = &server.Handler{Run: func(genericPeer cellnet.GenericPeer, event cellnet.Event) error {
	// 获取聊天消息内容
	chatMsg := event.Message().(*message.Message)
	// 无论客户端是否在线，都需要先保存消息存储库
	responseMsg := message.Message{
		From:       chatMsg.From,
		To:         chatMsg.To,
		Cmd:        message.CommandType_CHAT_RESP,
		CreateTime: time.Now().Unix() / 1e6,
		MsgType:    chatMsg.MsgType,
		ChatType:   chatMsg.ChatType,
		GroupId:    chatMsg.GroupId,
		Content:    chatMsg.Content,
		Extras:     chatMsg.Extras,
	}
	// 暂时选择忽略保存数据库出现的错误
	_ = server.SaveMsgRepo(&responseMsg)
	// 判断消息类型
	if chatMsg.ChatType == message.ChatType_PRIVATE { // 私聊
		// 获取好友ID对应的SessionID
		friendIdString := strconv.FormatInt(chatMsg.To, 10)
		sessionIdString, exist, err := server.Store.Find(friendIdString)
		if !exist || err != nil {
			// 如果用户当前不在线，需要保存消息到同步库
			_ = server.SavePrivateMsgSyncRepo(&responseMsg)
			return err
		}
		var sessionId int64
		if err := utils.StrToInt(sessionIdString, &sessionId); err != nil {
			return err
		}
		// 获取session，并且给对应的人发送消息
		session := genericPeer.(cellnet.SessionAccessor).GetSession(sessionId)
		if session == nil {
			// 因为异常原因掉线，需要保存消息到同步库
			_ = server.SavePrivateMsgSyncRepo(&responseMsg)
		} else {
			// 客户端在线，直接推送消息
			session.Send(&responseMsg)
		}
		return nil
	} else if chatMsg.ChatType == message.ChatType_PUBLIC { // 群聊
		// 群组用户ID列表（不包括自己）
		toIds := []int64{1, 2, 3}
		for _, toId := range toIds {
			// 获取成员ID对应的SessionID
			memberIdString := strconv.FormatInt(toId, 10)
			sessionIdString, exist, err := server.Store.Find(memberIdString)
			if !exist || err != nil {
				// 如果用户当前不在线，需要保存消息到同步库
				_ = server.SavePublicMsgSyncRepo(toId, &responseMsg)
				return err
			}
			var sessionId int64
			if err := utils.StrToInt(sessionIdString, &sessionId); err != nil {
				return err
			}
			// 获取session，并且给对应的人发送消息
			session := genericPeer.(cellnet.SessionAccessor).GetSession(sessionId)
			if session == nil {
				// 因为异常原因掉线，需要保存消息到同步库
				_ = server.SavePublicMsgSyncRepo(toId, &responseMsg)
			} else {
				// 客户端在线，直接推送消息
				session.Send(&responseMsg)
			}
		}
	}
	return nil
}}
