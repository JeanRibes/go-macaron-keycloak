package web

import (
	"context"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
	"log"
)

var Context = context.TODO()

func CreateServer() *macaron.Macaron {
	m := macaron.Classic()
	m.Use(macaron.Renderer()) //pour le templates HTML
	m.Use(session.Sessioner(session.Options{
		Provider:       "file",
		ProviderConfig: "local-sessions-cache",
	})) // pour avoir des données de session comme PHP

	m.Use(csrf.Csrfer())
	return m
}

func handleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
