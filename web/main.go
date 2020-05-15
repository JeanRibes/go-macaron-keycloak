package web

import (
	"context"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/pongo2"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
	"log"
)

var Context = context.TODO()

func CreateServer() *macaron.Macaron {
	m := macaron.Classic()
	m.Use(session.Sessioner(session.Options{
		Provider:       "file",
		ProviderConfig: "local-sessions-cache",
	})) // pour avoir des données de session comme PHP
	m.Use(csrf.Csrfer())
	m.Use(pongo2.Pongoer(pongo2.Options{Directory: "templates"})) //pour le templates HTML

	m.Use(func(ctx *macaron.Context) {
		ctx.Header().Set("X-BdE", "Si tu vois ceci, il faut que tu rejoignes l'équipe SIA si tu n'y es pas déjà :)")
		ctx.Header().Set("X-Powered-By", "go-macaron")
	})
	return m
}

func handleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
