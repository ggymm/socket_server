package handle

import (
	"github.com/davyxu/cellnet"
	"socket_server/server"
)

var UserHandle = &server.Handler{Run: func(genericPeer cellnet.GenericPeer, event cellnet.Event) error {

	return nil
}}
