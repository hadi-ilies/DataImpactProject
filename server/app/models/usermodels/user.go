package usermodels

//User: User model as given inside the given "dataset"
type User struct {
	Id        string      `json:"id" bson:"_id"`
	Password  string      `json:"password" bson:"password"`
	IsActive  bool        `json:"isActive" bson:"isActive"`
	Balance   string      `json:"balance" bson:"balance"`
	Age       interface{} `json:"age" bson:"age"` //sometimes the age is a string and sometimes it's an int inside the given dataset (I don't know if it's normal) I put an interface type in case it is
	Name      string      `json:"name" bson:"name"`
	Gender    string      `json:"gender" bson:"gender"`
	Company   string      `json:"company" bson:"company"`
	Email     string      `json:"email" bson:"email"`
	Phone     string      `json:"phone" bson:"phone"`
	Address   string      `json:"address" bson:"address"`
	About     string      `json:"about" bson:"about"`
	Registred string      `json:"registred" bson:"registred"`
	Latitude  float64     `json:"latitude" bson:"latitude"`
	Longitude float64     `json:"longitude" bson:"longitude"`
	Tags      []string    `json:"tags" bson:"tags"`
	Friends   []Friend    `json:"friends" bson:"friends"`
	Data      string      `json:"data" bson:"data"`

	//not asked
	Token string `json:"-" bson:"token"`
}
