package handle

import "github.com/davyxu/golog"

var LogHandle = golog.New("server.handle")

const (
	LoginSuccessLog = "用户登陆成功，ID为："
)
