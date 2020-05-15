package main

import (
	"evigo/db"
	"evigo/web"
	"github.com/go-macaron/pongo2"
	"gopkg.in/macaron.v1"
)

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	t := macaron.Classic()
	t.Use(pongo2.Pongoer(pongo2.Options{
		Directory: "pongot",
	}))

	t.Get("/", func(ctx *macaron.Context) {
		ctx.HTML(200, "test")
	})
	t.Run()

	m := web.CreateServer() //  https://go-macaron.com/
	println("serv created")

	db.Connect()
	defer db.Disconnect()

	web.SetupRemoteAuth()
	web.SetupRoutes(m)

	//log.Fatal(http.ListenAndServe("0.0.0.0:8000", m))
	m.Run()
}
