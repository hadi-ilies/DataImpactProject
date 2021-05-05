package usermodels

//Friend: User can have an array of friend
type Friend struct {
	Id   uint64 `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}
