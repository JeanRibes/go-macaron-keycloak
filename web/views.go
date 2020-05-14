package web

import (
	"ami/db"
	"fmt"
	"github.com/go-macaron/csrf"
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

func TableauAdherents(ctx *macaron.Context, x csrf.CSRF) {
	valides := db.ListValidAdherents()
	invalides := db.ListUnpaidAdherents()

	ctx.Data["Valides"] = valides
	ctx.Data["Invalides"] = invalides

	ctx.Data["csrf_token"] = x.GetToken()

	ctx.HTML(200, "gestion/liste")
}

func ModifPaiement(ctx *macaron.Context) {
	form := ctx.Req.Form
	fmt.Printf("data: %s\n %v\n", form, form)
	modifs := []string{}
	for name := range ctx.Req.Form {
		if name == "_csrf" {
			continue
		}
		modifs = append(modifs, name)
	}
	db.UpdateMultiplePayments(modifs)
	ctx.JSON(200, modifs)
}

func RechercheUsersKC(ctx *macaron.Context, x csrf.CSRF) {
	ctx.Data["csrf_token"] = x.GetToken()
	ctx.HTML(200, "gestion/recherche")
}
func RechercheUsersKcRresult(ctx *macaron.Context, x csrf.CSRF) {
	ctx.Data["csrf_token"] = x.GetToken()

	email := ctx.Req.Form["email"]
	ctx.Data["Resultats"] = true
	ctx.Data["Users"] = SearchUsersByEmail(email[0])
	ctx.HTML(200, "gestion/recherche")
}

func AjoutUser(ctx *macaron.Context) {
	println(ctx.Req.Form["username"][0])
	println(ctx.Req.Form["id"][0])
	user := GetUserById(ctx.Req.Form["id"][0])
	adherent := initAdherent()
	adherent.UserToAdherent(user)

	db.CreateAdherent(adherent)

	l := "/adherent/" + adherent.Username + "/"
	println(l)
	fmt.Printf("%s\n", l)
	ctx.Redirect(l)
}

func ViewAdherent(ctx *macaron.Context) {
	ad, err := db.FindAdherentByUsername(ctx.Params(":username"))
	showError(ctx, err)
	ctx.Data["Adherent"] = ad
	ctx.HTML(200, "profile")
}
