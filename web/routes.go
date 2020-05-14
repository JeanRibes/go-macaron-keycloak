package web

import (
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
	"net/http"
)

func SetupRoutes(m *macaron.Macaron) {
	m.Get("/", Index)

	//authentification Keycloak
	m.Get("/start/", OidcStart)
	m.Get("/return/", OidcFinish)
	m.Get("/profile", AuthenticationMiddleware, MyProfile)

	m.Get("/list", BureauMiddleware, TableauAdherents)
	m.Post("/modifpay", csrf.Validate, BureauMiddleware, ModifPaiement)

	m.Get("/test", func(w http.ResponseWriter, r *http.Request, ctx *macaron.Context) {
		ctx.Header().Set("Content-Type", "text/html")
		ctx.RawData(200, []byte("alkdfsj sldf  kl klj jkl jkljl kjlk jk"))
	})

	// flash sert apparemmetn à stocker des données uniquement pour la requête suivant
	m.Get("/flash", func(ctx *macaron.Context, f *session.Flash) {
		f.Set("A", "aaaa")
		f.Success("yes!!!", true)
		f.Error("opps...", true)
		f.Info("aha?!", true)
		f.Warning("Just be careful.", true)
		ctx.HTML(200, "flash")
	})
}
