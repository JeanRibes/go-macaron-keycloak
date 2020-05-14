package db

import (
	"time"
)

type IdToken struct {
	AuthToken    string
	RefreshToken string
}
type Adherent struct { //les champs CamelCase sont convertis en lowercase
	//ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	//récup de keycloak
	Username  string  ``
	Email     string  ``
	Auth      IdToken ``
	FirstName string  ``
	LastName  string
	Gender    string `` //c'est un attribut sur KC

	JoinedAt     time.Time ``
	Commentaires string    //pour les admins
	APaye        bool
}
