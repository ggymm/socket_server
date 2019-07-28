package server

import (
	"encoding/json"
	"github.com/davyxu/cellnet"
	"socket_server/server/message"
	"socket_server/utils"
	"socket_server/web/models"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	Run func(genericPeer cellnet.GenericPeer, event cellnet.Event) error
}

func (handler *Handler) Handle(genericPeer cellnet.GenericPeer, event cellnet.Event) error {
	return handler.Run(genericPeer, event)
}

type handler interface {
	Handle(genericPeer cellnet.GenericPeer, event cellnet.Event) error
}

var Handlers = make(map[message.CommandType]handler)

func SaveMsgRepo(msg *message.Message) error {
	// 保存消息到存储库，持久化保存。用于读漫游消息
	// 使用mysql存储消息内容
	messageInfo := new(models.Message)
	messageInfo.FromId = uint64(msg.From)
	messageInfo.ToId = uint64(msg.To)
	messageInfo.CreateTime = time.Unix(0, msg.CreateTime*1e6)
	messageInfo.MsgType = uint32(msg.GetMsgType())
	messageInfo.ChatType = uint32(msg.GetChatType())
	messageInfo.GroupId = msg.GroupId
	messageInfo.Content = msg.Content
	messageInfo.Extras = msg.Extras
	return models.SaveMessage(messageInfo)
}

func SavePrivateMsgSyncRepo(msg *message.Message) error {
	// 保存消息到同步仓库，用于离线消息临时存储
	// 使用redis作为同步仓库
	// PRIVATE_FromId_ToId
	msgKey := strings.Join(
		[]string{message.ChatType_PRIVATE.String(),
			strconv.FormatInt(msg.From, 10),     // 发送者ID
			strconv.FormatInt(msg.To, 10)}, "_") // 接收者ID
	if err := Store.EnQueue(msgKey, msg.String()); err != nil {
		return err
	}
	return nil
}

func SavePublicMsgSyncRepo(toId int64, msg *message.Message) error {
	// 保存消息到同步仓库，用于离线消息临时存储
	// 使用redis作为同步仓库
	// PUBLIC_GroupId_ToId
	msgKey := strings.Join(
		[]string{message.ChatType_PUBLIC.String(),
			msg.GroupId,                       // 群组ID（最好不要用toID，会引起歧义）
			strconv.FormatInt(toId, 10)}, "_") // 群成员ID（除了自己）
	if err := Store.EnQueue(msgKey, msg.String()); err != nil {
		return err
	}
	return nil
}

func GetMsgSyncCount(isGroup bool, friendIds []int64, msg *message.Message) ([]map[string]string, error) {
	chatUnReadMsgList := make([]map[string]string, 0)
	for _, friendId := range friendIds {
		friendIdString := strconv.FormatInt(friendId, 10)
		msgKey := strings.Join(
			[]string{utils.If(isGroup, message.ChatType_PRIVATE.String(), message.ChatType_PUBLIC.String()),
				friendIdString,                        // 包括好友ID和群组ID
				strconv.FormatInt(msg.From, 10)}, "_") // 自己ID
		count, err := Store.CountQueued(msgKey)
		if err != nil {
			return chatUnReadMsgList, err
		}
		chatUnReadMsgList = append(chatUnReadMsgList, map[string]string{
			"ID":          friendIdString,
			"UnReadCount": strconv.FormatInt(int64(count), 10),
		})
	}
	return chatUnReadMsgList, nil
}

func GetMsgSync(msg *message.Message) ([]string, error) {
	var contentMap map[string]string
	if err := json.Unmarshal([]byte(msg.Content), &contentMap); err == nil {
		return nil, err
	}
	chartType := contentMap["chartType"]
	friendIdString := contentMap["friendId"]
	isGroup := strings.EqualFold(chartType, message.ChatType_PUBLIC.String())
	msgKey := strings.Join(
		[]string{utils.If(isGroup, message.ChatType_PRIVATE.String(), message.ChatType_PUBLIC.String()),
			friendIdString,                        // 包括好友ID和群组ID
			strconv.FormatInt(msg.From, 10)}, "_") // 自己ID
	return Store.GetQueue(msgKey)
}
