package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/coreos/go-oidc"
	"github.com/go-macaron/session"
	"golang.org/x/oauth2"
	"gopkg.in/macaron.v1"
	"net/http"
)

var contexte = context.TODO()

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
func showError(ctx *macaron.Context, err error) {
	if err != nil {
		ctx.Data["Err"] = err
		ctx.HTML(500, "error")
	}
}

var verifier *oidc.IDTokenVerifier

func main() {
	m := macaron.Classic() //  https://go-macaron.com/
	m.Use(macaron.Renderer())
	m.Use(session.Sessioner())
	//m.Use(macaron.Renderer(macaron.RenderOptions{Directory: "./templates"}))

	keycloak, kerr := oidc.NewProvider(contexte, "http://192.168.0.25/auth/realms/asso-insa-lyon")
	handleError(kerr)

	oauth2Config := oauth2.Config{
		ClientID:     "goami",
		ClientSecret: "229b12d3-c523-4546-814a-c6199f5379c4",
		RedirectURL:  "http://localhost:4000/ok/",

		// Discovery returns the OAuth2 endpoints.
		Endpoint: keycloak.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}
	verifier = keycloak.Verifier(&oidc.Config{ClientID: "goami"})

	m.Get("/", func(ctx *macaron.Context) {
		ctx.Data["Variable"] = "valeur"
		ctx.HTML(200, "test")
	})

	m.Get("/test", func(w http.ResponseWriter, r *http.Request, ctx *macaron.Context) {
		ctx.Header().Set("Content-Type", "text/html")
		ctx.RawData(200, []byte("alkdfsj sldf  kl klj jkl jkljl kjlk jk"))
	})

	// http://localhost:4000/ok/?state=se-souvenir-de-cette-string&session_state=7187b971-1077-46b2-8921-4e32028886b6&code=0664848a-1ad9-4acc-978d-ad66ba17a25f.7187b971-1077-46b2-8921-4e32028886b6.fbcc0084-f1ab-4088-aa87-2fa284e824ee
	m.Get("/oidc", func(w http.ResponseWriter, r *http.Request, ctx *macaron.Context, sess session.Store) {
		sess.Set("oidc-state", "se-souvenir-de-cette-string")
		http.Redirect(w, r, oauth2Config.AuthCodeURL("se-souvenir-de-cette-string"), http.StatusFound)
	})
	m.Get("/ok/", func(w http.ResponseWriter, r *http.Request, ctx *macaron.Context, sess session.Store) {
		// Verify state and errors.
		if ctx.Query("state") != "se-souvenir-de-cette-string" {
			showError(ctx, errors.New("mismatched state : quelqu'un a essayé de trifouiller dans la communication entre les serveurs !"))
		}

		oauth2Token, err := oauth2Config.Exchange(contexte, r.URL.Query().Get("code"))
		handleError(err)
		// Extract the ID Token from OAuth2 token.
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			// handle missing token
			showError(ctx, errors.New("missing token, verification failed"))
		}
		// Parse and verify ID Token payload.
		idToken, err := verifier.Verify(contexte, rawIDToken)
		if err != nil {
			// handle error
			showError(ctx, err)
		}

		// Extract custom claims
		var claims struct {
			Email    string `json:"email"`
			Verified bool   `json:"email_verified"`
		}
		if err := idToken.Claims(&claims); err != nil {
			// handle error
			showError(ctx, err)
		}
		fmt.Printf("auth: %s\n user_id: %s\n \n ath=%s", rawIDToken, idToken.Subject, idToken.AccessTokenHash)
		ctx.Data["Email"] = claims.Email
		ctx.Data["EmailVerified"] = claims.Verified
		ctx.Data["Keto"] = rawIDToken
		ctx.Data["OAuth2token"] = oauth2Token.AccessToken
		ctx.Data["OAuth2refresh"] = oauth2Token.RefreshToken
		ctx.Data["OAuth2type"] = oauth2Token.TokenType

		handleError(sess.Set("keycloak_token", oauth2Token.AccessToken))
		handleError(sess.Set("keycloak_refresh_token", oauth2Token.RefreshToken))
		handleError(sess.Set("email", claims.Email))
		handleError(sess.Set("username", idToken.Subject))
		ctx.HTML(200, "keycloak")
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

	//log.Fatal(http.ListenAndServe("0.0.0.0:8000", m))
	m.Run()
}
