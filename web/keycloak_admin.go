package web

import (
	"github.com/Nerzal/gocloak"
)

var client gocloak.GoCloak
var token *gocloak.JWT

func SetupAdminClient() {
	client = gocloak.NewClient(KeycloakUrl)
	var err error
	token, err = client.LoginClient(ClientId, ClientSecret, Realm)
	handleError(err)
	//users, err := client.GetUsers(token.AccessToken, "asso-insa-lyon", gocloak.GetUsersParams{})
	//users, err := client.GetUsers(token.AccessToken, "asso-insa-lyon", gocloak.GetUsersParams{Email: ""})
	//démontre que c'est interdit par la permission view-users
	/*user,err:=client.GetUserByID(token.AccessToken, "asso-insa-lyon", "6e22bf51-b91c-4b9a-a677-fb52689aff00")
	handleError(err)
	fmt.Printf("%s %s - %s %s\n", user.Email, user.Username, user.FirstName, user.LastName)
	handleError(err)*/

}

func SearchUsersByEmail(email string) *[]gocloak.User {
	println(email)
	users, err := client.GetUsers(token.AccessToken, Realm, gocloak.GetUsersParams{Email: email})
	handleError(err)
	return users
}

func GetUserById(id string) *gocloak.User {
	// &{6e22bf51-b91c-4b9a-a677-fb52689aff00 %!s(int64=1588971800117) dv0de7sdlphjqri8nnvrlayyy4q %!s(bool=true) %!s(bool=false) %!s(bool=true) jean ribes jean.christophe.ribes@sfr.fr  map[adhesionUserId:[24] birthday:[2020-05-19] category:[student] department:[FIMI] gender:[M] locale:[fr] memberId:[32] school:[INSA] study_year:[1A] terms_and_conditions:[1588971845]] [] [] map[impersonate:%!s(bool=false) manage:%!s(bool=false) manageGroupMembership:%!s(bool=false) mapRoles:%!s(bool=false) view:%!s(bool=true)]}
	user, err := client.GetUserByID(token.AccessToken, Realm, id)
	handleError(err)
	return user
}
