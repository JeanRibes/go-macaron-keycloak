package web

import (
	"ami/db"
	"fmt"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
)

func showError(ctx *macaron.Context, err error) {
	if err != nil {
		ctx.Data["Err"] = err
		ctx.HTML(500, "error")
	}
}

func Index(ctx *macaron.Context) {
	ctx.Data["Variable"] = "valeur"
	ctx.HTML(200, "test")
}

func MyProfile(ctx *macaron.Context, sess session.Store) {
	adherent, err := db.FindAdherentByUsername(fmt.Sprintf("%s", sess.Get("username")))
	showError(ctx, err)
	ctx.Data["Adherent"] = *adherent
	ctx.HTML(200, "profile")
}
