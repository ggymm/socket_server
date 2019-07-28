package server

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/gorillaws"
	"github.com/davyxu/cellnet/proc"
	_ "github.com/davyxu/cellnet/proc/gorillaws"
	"github.com/davyxu/golog"
	"socket_server/config"
	"socket_server/server/message"
)

var log = golog.New("server.start")

func StartServer() {
	// 获取配置文件
	peerType := config.Config.Get("cellnet.peerType").(string)
	name := config.Config.Get("cellnet.name").(string)
	addr := config.Config.Get("cellnet.addr").(string)
	procName := config.Config.Get("cellnet.procName").(string)
	// _ = golog.SetLevelByString(".", "info")
	// 创建事件处理队列
	eventQueue := cellnet.NewEventQueue()
	// 创建监听器，将事件传送给事件处理队列
	genericPeer := peer.NewGenericPeer(peerType, name, addr, eventQueue)
	// 设置接收消息模式，并且接收消息，处理消息
	proc.BindProcessorHandler(genericPeer, procName, func(event cellnet.Event) {
		switch msg := event.Message().(type) {
		// 连接成功
		case *cellnet.SessionAccepted:
			// 连接成功暂不做任何操作，只是允许连接
			log.Infoln("会话创建成功，会话ID：【", event.Session().ID(), "】")
		// 断开连接
		case *cellnet.SessionClosed:
			log.Infoln("会话关闭成功，会话ID：【", event.Session().ID(), "】")
		// 消息处理
		case *message.Message:
			// 此处会根据消息类型自动调用不同的处理方法
			if err := Handlers[msg.Cmd].Handle(genericPeer, event); err != nil {
				log.Errorln("会话出现错误，会话ID：【", event.Session().ID(), "】。", "消息内容：【", msg, "】")
			}
		}
	})
	// 启动侦听
	genericPeer.Start()
	// 事件队列开始循环
	eventQueue.StartLoop()
	// 等待退出消息，退出事件队列
	eventQueue.Wait()
}
