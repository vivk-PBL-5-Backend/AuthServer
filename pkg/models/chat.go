package models

type Chat struct {
	Username   string   `json:"username" bson:"_id"`
	Companions []string `json:"companions" bson:"companions"`
}
