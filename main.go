package main

import (
	"ami/db"
	"ami/web"
)

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	m := web.CreateServer() //  https://go-macaron.com/

	db.Connect()
	defer db.Disconnect()

	web.SetupRemoteAuth()
	web.SetupRoutes(m)

	//log.Fatal(http.ListenAndServe("0.0.0.0:8000", m))
	m.Run()
}
