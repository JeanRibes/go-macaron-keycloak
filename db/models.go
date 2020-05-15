package db

import (
	"github.com/Nerzal/gocloak"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IdToken struct {
	AuthToken    string
	RefreshToken string
}
type Adherent struct { //les champs CamelCase sont convertis en lowercase
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	//récup de keycloak
	Username  string  ``
	Email     string  ``
	Auth      IdToken ``
	FirstName string  ``
	LastName  string
	Gender    string `` //c'est un attribut sur KC

	JoinedAt     time.Time `bson:"joined_at"`
	Commentaires string    //pour les admins
	APaye        bool      `bson:"a_paye"`
	RoleBureau   bool      `bson:"role_bureau"`
}

func (this *Adherent) UserToAdherent(user *gocloak.User) {
	this.LastName = user.LastName
	this.FirstName = user.FirstName
	this.Email = user.Email
	this.Username = user.Username
	for key, values := range user.Attributes {
		if key == "gender" {
			this.Gender = values[0]
		}
	}
}
func (this *Adherent) DisplayGender() string {
	if this.Gender == "W" {
		return "femme"
	}
	if this.Gender == "M" {
		return "homme"
	} else {
		return "_"
	}
}

func (this *Adherent) String() string {
	return this.FirstName + " " + this.LastName
}
