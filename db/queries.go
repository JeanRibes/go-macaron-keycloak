package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func CreateAdherent(adherent *Adherent) {
	is, err := AdherentsCollection.InsertOne(Context, adherent)
	handleError(err)
	handleError(AdherentsCollection.FindOne(Context, bson.D{{"_id", is.InsertedID}}).Decode(&adherent)) //vérif
}

func UpdateAdherent(adherent *Adherent) *Adherent {
	filter := bson.D{
		{"username", adherent.Username},
		{"_id", adherent.ID},
	}
	update := bson.D{{"$set", adherent}}
	_, err := AdherentsCollection.UpdateOne(Context, filter, update)
	handleError(err)
	return adherent
}

// Renvoie tous les adhérents qui ont payé
func ListValidAdherents() []Adherent {
	var valides []Adherent
	cur, err := AdherentsCollection.Find(Context, bson.D{{"a_paye", true}}, options.Find())
	handleError(err)
	handleError(cur.All(Context, &valides))

	return valides
}

func ListUnpaidAdherents() []Adherent {
	var non_payes []Adherent
	cur, err := AdherentsCollection.Find(Context, bson.D{{"a_paye", false}}, options.Find())
	handleError(err)
	handleError(cur.All(Context, &non_payes))

	return non_payes
}

func UpdateMultiplePayments(usernames []string) {

	filter := bson.D{}
	for _, username := range usernames {
		filter = append(filter, bson.E{Key: "username", Value: username})
	}
	update := bson.D{{"$set", bson.D{
		{"a_paye", true},
	}}}
	_, err := AdherentsCollection.UpdateMany(Context, filter, update)
	handleError(err)
}
