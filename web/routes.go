package web

import (
	"github.com/go-macaron/csrf"
	"gopkg.in/macaron.v1"
)

func SetupRoutes(m *macaron.Macaron) {
	m.Get("/", Index)

	//authentification Keycloak
	m.Get("/start/", OidcStart)
	m.Get("/return/", OidcFinish)
	m.Get("/profile", AuthenticationMiddleware, MyProfile)

	m.Get("/adherents/", BureauMiddleware, TableauAdherents)
	m.Get("/adherents/:username/", BureauMiddleware, ViewAdherent)
	m.Post("/modifpay/", csrf.Validate, BureauMiddleware, ModifPaiement)
	m.Get("/keycloak/recherche/", BureauMiddleware, RechercheUsersKC)
	m.Post("/keycloak/recherche/", csrf.Validate, BureauMiddleware, RechercheUsersKcRresult)
	m.Post("/keycloak/ajout/", csrf.Validate, BureauMiddleware, AjoutUser)
}
