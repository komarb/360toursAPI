package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Link struct {
	ID string `json:"id"`
	F  int    `json:"f"`
	O  int    `json:"o"`
}
type Picture struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Image string `json:"image"`
	Links []Link `json:"links"`
}
type Tour struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	MainID   string             `json:"mainId" bson:"mainId"`
	Pictures []Picture          `json:"pictures"`
}
