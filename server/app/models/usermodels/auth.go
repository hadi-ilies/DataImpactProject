package usermodels

//The email is more accurate and common than an ID therefore I decided to replace the auth id + password by email + password
type Auth struct {
	Email    string `json:"email"  bson:"email"`
	Password string `json:"password" bson:"password"`
}
