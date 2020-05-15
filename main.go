package main

import (
	"evigo/db"
	"evigo/web"
)

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	m := web.CreateServer() //  https://go-macaron.com/

	web.SetupRemoteAuth() //initialise la connexion à Keycloak

	db.Connect()
	defer db.Disconnect()

	web.SetupRoutes(m)

	//log.Fatal(http.ListenAndServe("0.0.0.0:8000", m))
	m.Run()
}
