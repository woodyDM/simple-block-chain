package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/woodyDM/gframe/public/gf"
)

var profile = flag.String("p", "", "set profile name.")

func main() {
	flag.Parse()
	gf.InitConfig(*profile)
	bindPort := fmt.Sprintf(":%d", gf.Conf.Port)

	app, err := gf.CreateMyServer(bindPort)
	if err != nil {
		panic(err)
	}
	rg := app.Engine.Group("/api")
	registerRouters(rg)
	gf.Log.Printf("Starting server at port [%s]", bindPort)
	app.Start()
	app.Await()
}

/**
to register routers here
*/
func registerRouters(r *gin.RouterGroup) {

	r.Use(gf.DefaultCookieHandler)
	r.Use(gf.DefaultUserAgentWhiteListHandler)

	r.GET("/header", gf.Route(func(ctx *gin.Context) (interface{}, error) {
		return ctx.HandlerNames(), nil
	}))

}
