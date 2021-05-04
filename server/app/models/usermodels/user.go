package usermodels

import (
	//We will be using mongo that’s why we use a bson.ObjectId
	"time"
)

//User: User model as given inside the given "dataset"
type User struct {
	Id        string    `json:"id" bson:"_id"`
	Password  string    `json:"password" bson:"password"`
	IsActive  bool      `json:"isActive" bson:"isActive"`
	Balance   string    `json:"balance" bson:"balance"`
	Age       uint8     `json:"age" bson:"age"`
	Name      string    `json:"name" bson:"name"`
	Gender    string    `json:"gender" bson:"gender"`
	Company   string    `json:"company" bson:"company"`
	Email     string    `json:"email" bson:"email"`
	Phone     string    `json:"phone" bson:"phone"`
	Address   string    `json:"address" bson:"address"`
	About     string    `json:"about" bson:"about"`
	Registred time.Time `json:"registred" bson:"registred"`
	Latitude  float64   `json:"latitude" bson:"latitude"`
	Longitude float64   `json:"longitude" bson:"longitude"`
	Tags      []string  `json:"tags" bson:"tags"`
	Friends   []Friend  `json:"friends" bson:"friends"`
	Data      string    `json:"data" bson:"data"`

	//not asked
	Token string `json:"-" bson:"token"`
}
