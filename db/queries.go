package db

import (
	"go.mongodb.org/mongo-driver/bson"
)

func FindAdherentByUsername(username string) (*Adherent, error) {
	var adherent *Adherent
	sr := AdherentsCollection.FindOne(Context, bson.D{
		{"username", username},
	})
	if sr.Err() != nil {
		return nil, sr.Err() //erreur non trouvée ...
	}

	handleError(sr.Decode(&adherent)) //erreur spéciale
	return adherent, nil
}

func CreateAdherent(adherent Adherent) {
	is, err := AdherentsCollection.InsertOne(Context, adherent)
	handleError(err)
	handleError(AdherentsCollection.FindOne(Context, bson.D{{"_id", is.InsertedID}}).Decode(&adherent)) //vérif
}

func UpdateAdherent(adherent *Adherent) *Adherent {
	//AdherentsCollection.UpdateOne(Context, )
	return adherent
}
