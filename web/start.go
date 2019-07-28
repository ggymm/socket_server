package web

import (
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/router"
	"github.com/kataras/iris/middleware/logger"
	"socket_server/config"
	"socket_server/constant"
	"socket_server/web/controllers"
	"socket_server/web/database"
)

func newApp() (api *iris.Application) {
	api = iris.New()
	api.Use(logger.New())

	api.OnErrorCode(iris.StatusNotFound, func(context iris.Context) {
		_, _ = context.JSON(controllers.ApiResource(iris.StatusNotFound, nil, constant.StatusNotFound))
	})
	api.OnErrorCode(iris.StatusInternalServerError, func(context iris.Context) {
		_, _ = context.WriteString(constant.StatusInternalServerError)
	})

	iris.RegisterOnInterrupt(func() {
		_ = database.DB.Close()
	})

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, //允许通过的主机名称
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	v1 := api.Party("/api/v1", crs).AllowMethods(iris.MethodOptions)
	{
		v1.PartyFunc("/im", func(imApi router.Party) {
		})
	}

	return
}

func StartWebServer() {
	app := newApp()
	addr := config.Config.Get("web.addr").(string)
	_ = app.Run(iris.Addr(addr))
}
