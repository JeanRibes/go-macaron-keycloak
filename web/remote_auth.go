package web

import (
	"ami/db"
	"errors"
	"fmt"
	"github.com/coreos/go-oidc"
	"github.com/go-macaron/session"
	"golang.org/x/oauth2"
	"gopkg.in/macaron.v1"
	"net/http"
	"os"
	"time"
)

var oauth2Config oauth2.Config
var verifier *oidc.IDTokenVerifier

var KeycloakUrl string
var Realm string
var ClientId string
var ClientSecret string

type Claims struct {
	Email     string   `json:"email"`
	Verified  bool     `json:"email_verified"`
	LastName  string   `json:"family_name"`
	FirstName string   `json:"given_name"`
	Gender    string   `json:"gender"`
	Username  string   `json:"preferred_username"`
	AmiRoles  []string `json:"ami_roles"`
}

func (claims Claims) IsBureau() bool {
	for _, role := range claims.AmiRoles {
		if role == "bureau" {
			return true
		}
	}
	return false
}

func SetupRemoteAuth() {
	KeycloakUrl = os.Getenv("KEYCLOAK_URL")
	if KeycloakUrl == "" {
		KeycloakUrl = "http://192.168.0.25"
	}
	Realm = os.Getenv("KEYCLOAK_REALM")
	if Realm == "" {
		Realm = "asso-insa-lyon"
	}
	ClientId = os.Getenv("CLIENT_ID")
	if ClientId == "" {
		ClientId = "goami"
	}
	ClientSecret = os.Getenv("CLIENT_SECRET")
	if ClientSecret == "" {
		ClientSecret = "229b12d3-c523-4546-814a-c6199f5379c4"
	}

	keycloak, kerr := oidc.NewProvider(Context, KeycloakUrl+"/auth/realms/"+Realm)
	handleError(kerr)

	oauth2Config = oauth2.Config{
		ClientID:     ClientId,
		ClientSecret: ClientSecret,
		RedirectURL:  "http://localhost:4000/return/",

		// Discovery returns the OAuth2 endpoints.
		Endpoint: keycloak.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}
	verifier = keycloak.Verifier(&oidc.Config{ClientID: "goami"})

	SetupAdminClient()
}

func OidcStart(w http.ResponseWriter, r *http.Request, sess session.Store) {
	handleError(sess.Set("oidc-state", "se-souvenir-de-cette-string"))
	http.Redirect(w, r, oauth2Config.AuthCodeURL("se-souvenir-de-cette-string"), http.StatusFound)
}

func OidcFinish(r *http.Request, ctx *macaron.Context, sess session.Store) {
	// Verify state and errors.
	if ctx.Query("state") != "se-souvenir-de-cette-string" {
		showError(ctx, errors.New("mismatched state : quelqu'un a essayé de trifouiller dans la communication entre les serveurs !"))
	}

	oauth2Token, err := oauth2Config.Exchange(Context, r.URL.Query().Get("code"))
	handleError(err)
	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok { // handle missing token
		showError(ctx, errors.New("missing token, verification failed"))
	}

	idToken, err := verifier.Verify(Context, rawIDToken) // Parse and verify ID Token payload.
	if err != nil {
		showError(ctx, err)
	}

	// Extract custom claims
	// {"exp":1589443715,"iat":1589443415,"auth_time":1589443415,"jti":"eab64e94-9681-419f-9b9e-ff3414309731","iss":"http://192.168.0.25/auth/realms/asso-insa-lyon","aud":"goami","sub":"6e22bf51-b91c-4b9a-a677-fb52689aff00","typ":"ID","azp":"goami","session_state":"8f08ed07-e491-4054-811c-f8c4eb143509","acr":"1","email_verified":true,"gender":"W","name":"jean ribes","preferred_username":"dv0de7sdlphjqri8nnvrlayyy4q","given_name":"jean","locale":"fr","family_name":"ribes","email":"jean.christophe.ribes@sfr.fr"}
	var claims Claims
	if err := idToken.Claims(&claims); err != nil {
		showError(ctx, err)
	}
	if !claims.Verified { // si l'email n'est pas vérifié
		showError(ctx, errors.New(fmt.Sprintf("Vous devez vérifier votre e-mail ()", claims.Email)))
		return
	}

	ctx.Data["Email"] = claims.Email
	ctx.Data["EmailVerified"] = claims.Verified
	ctx.Data["Keto"] = rawIDToken
	ctx.Data["OAuth2token"] = oauth2Token.AccessToken
	ctx.Data["OAuth2refresh"] = oauth2Token.RefreshToken
	ctx.Data["OAuth2type"] = oauth2Token.TokenType
	ctx.Data["Username"] = claims.Username

	handleError(sess.Set("keycloak_token", oauth2Token.AccessToken))
	handleError(sess.Set("keycloak_refresh_token", oauth2Token.RefreshToken))
	if claims.IsBureau() {
		handleError(sess.Set("role_bureau", true))
	}
	//handleError(sess.Set("email", claims.Email))
	handleError(sess.Set("username", claims.Username))
	adherent, errnotfound := db.FindAdherentByUsername(claims.Username)
	if errnotfound != nil {
		newAdherent := claimsToAdherent(claims)
		newAdherent.Auth = db.IdToken{
			AuthToken:    oauth2Token.AccessToken,
			RefreshToken: oauth2Token.RefreshToken,
		}
		db.CreateAdherent(newAdherent)
	} else {
		updateAdherentFromClaims(claims, adherent)
		updateAdherentFromToken(oauth2Token, adherent)
		db.UpdateAdherent(adherent)
	}

	ctx.HTML(200, "keycloak")
}

func updateAdherentFromClaims(claims Claims, adherent *db.Adherent) {
	adherent.Username = claims.Username
	adherent.FirstName = claims.FirstName
	adherent.LastName = claims.LastName
	adherent.Email = claims.Email
	adherent.Gender = claims.Gender
	adherent.RoleBureau = claims.IsBureau()
}
func initAdherent() *db.Adherent {
	return &db.Adherent{
		Commentaires: "ajouté par Keycloak",
		JoinedAt:     time.Now(),
		APaye:        false,
	}
}
func claimsToAdherent(claims Claims) *db.Adherent {
	adherent := initAdherent()
	updateAdherentFromClaims(claims, adherent)
	return adherent
}
func updateAdherentFromToken(token *oauth2.Token, adherent *db.Adherent) {
	adherent.Auth = db.IdToken{
		AuthToken:    token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
}

func AuthenticationMiddleware(ctx *macaron.Context, store session.Store) {
	keycloak_token := store.Get("keycloak_token")
	if keycloak_token == nil {
		showError(ctx, errors.New("Vous n'êtes pas authentifié !"))
	} else {
		// correct
		// ctx.Map(&adherent)
	}
}
func BureauMiddleware(ctx *macaron.Context, store session.Store) {
	bureau_role := store.Get("role_bureau")
	if bureau_role == nil {
		showError(ctx, errors.New("Vous n'êtes pas membre du Bureau ! permission non accordée"))
	}
}
