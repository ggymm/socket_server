package handle

import (
	"encoding/json"
	"github.com/davyxu/cellnet"
	"socket_server/constant"
	"socket_server/server"
	"socket_server/server/message"
	"strconv"
	"time"
)

var LoginHandle = &server.Handler{Run: func(genericPeer cellnet.GenericPeer, event cellnet.Event) error {
	// 获取登陆消息内容
	loginMsg := event.Message().(*message.Message)
	// 如果登陆成功，需要记录用户ID和会话ID
	userIdString := strconv.FormatInt(loginMsg.From, 10)
	sessionIdString := strconv.FormatInt(event.Session().ID(), 10)
	if err := server.Store.Save(userIdString, sessionIdString); err != nil {
		return err
	}
	// 获取聊天列表ID，同时获取未读消息数
	// TODO 这里应该是遍历聊天列表
	var chatUnReadMsgListString []byte
	if loginMsg.From == 2 {
		chatUnReadMsgList, getMessageErr := server.GetMsgSyncCount(false, []int64{1}, loginMsg)
		if getMessageErr != nil {
			return getMessageErr
		}
		// 此处的JSON转换库，一定不会出现问题
		chatUnReadMsgListString, _ = json.Marshal(chatUnReadMsgList)
	}
	LogHandle.Infoln(LoginSuccessLog, "【", loginMsg.From, "】")
	// 生成返回客户端登录成功的消息并且返回客户端
	event.Session().Send(&message.Message{
		From:       constant.ServerId,
		To:         loginMsg.From,
		Cmd:        message.CommandType_LOGIN_RESP,
		CreateTime: time.Now().Unix() / 1e6,
		MsgType:    message.MsgType_TEXT,
		ChatType:   message.ChatType_PRIVATE,
		GroupId:    constant.Empty,
		Content:    constant.LoginSuccess,
		Extras:     string(chatUnReadMsgListString),
	})
	return nil
}}
