package usercontrollers

import (
	"DataImpactProject/server/app/db"
	"DataImpactProject/server/app/models/usermodels"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

//Create: The JSON input file must be serialized and then saved to a MongoDB database. The JSON file will contain an array of users. You will then have to serialize the data concurrently and then insert it into the database without processing the already inserted entries again.The password must be encrypted with bcrypt and the hash of it inserted into the database.In addition to inserting this data into the database, you will have to generate one file per entry with the entry id as the filename, this file should contain only the Data field
func Create(ctx *gin.Context) {
	var users []usermodels.User
	err := ctx.Bind(&users)
	if err != nil {
		fmt.Println("Error = ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error while decoding",
		})
		return
	}

	session, err := db.GetMongoSession()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error while getting mongo session",
		})
		return
	}
	defer session.Close()
	c := session.DB("dataImpact").C("users")
	for _, user := range users {
		//TODO check id user exist with this email
		err = c.Find(bson.M{"email": user.Email}).One(&user)
		if err == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Error account with this email already exist",
			})
			continue
		}
		//encrypt password
		user.Password = hashAndSalt([]byte(user.Password))
		err = c.Insert(&user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Error while inserting data",
			})
			return
		}
	}
	ctx.JSON(http.StatusOK, users)
	// wait := new(sync.WaitGroup)

	// for i, user := range users {
	// 	wait.Add(1)
	// 	fmt.Println("user: ", i, " = ", user.Name)
	//     go func(j int) {
	// 		defer wait.Done()
	// 		deserializedObject, err := doDeserialization(user)
	// 		if err != nil {
	// 			// add error handling here
	// 		}
	// 		objs[j] = deserializedObject
	// 	}(j)
	// }
}

//Login: The user is able to connect to his profile in order to access all the data assigned to him. We consider that the user has an email and a password
func Login(ctx *gin.Context) {
	var (
		authUser usermodels.Auth
		user     usermodels.User
		err      error
	)

	err = ctx.Bind(&authUser)
	if err != nil {
		fmt.Println("Error = ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error while decoding",
		})
		return
	}
	session, err := db.GetMongoSession()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error while getting mongo session",
		})
		return
	}
	defer session.Close()
	c := session.DB("dataImpact").C("users")
	err = c.Find(bson.M{"email": authUser.Email}).One(&user)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Error account with this email doesn't exist",
		})
		return
	}
	//CHECK PASSWORD
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authUser.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Error: incorrect password has been inserted",
		})
		return
	}
	//account exist
	validToken, err := getJWT()
	fmt.Println(validToken)
	//Note: it is not really nessesary to save the token but can be cool to keep it if I want to avoid sending ID in request in order to get/delete/update my user
	user.Token = validToken
	//update
	err = c.Update(bson.M{"email": authUser.Email}, user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error: while saving token",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"token": user.Token,
	})
}

//Delete: delete a user with its Id, as well as the generated file (I'm doing as the subject want) but if we were in a real case. I would deleted the ID param from the request and I would allowed a user to delete ONLY his own account thx to his token
func Delete(ctx *gin.Context) {
	_, err := IsAuthorized(ctx.Request)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Error: Access token not found or not valid",
		})
		return
	}

	session, err := db.GetMongoSession()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error while getting mongo session",
		})
		return
	}
	defer session.Close()
	c := session.DB("dataImpact").C("users")
	userID := ctx.Param("id")
	err = c.Remove(bson.M{"_id": userID})
	if err != nil {
		//TODO check status code
		ctx.JSON(http.StatusNoContent, gin.H{
			"message": "User does not exist",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

//GetAllUsers: retrieve a list of users
func GetAllUsers(ctx *gin.Context) {
	_, err := IsAuthorized(ctx.Request)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Error: Access token not found or not valid",
		})
		return
	}

	session, err := db.GetMongoSession()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error while getting mongo session",
		})
		return
	}
	defer session.Close()
	c := session.DB("dataImpact").C("users")
	var users []usermodels.User
	err = c.Find(nil).All(&users)
	if err != nil {
		ctx.JSON(http.StatusNoContent, gin.H{
			"message": "Error: the DB is empty",
		})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

//GetUserByID: retrieve a user
func GetUserByID(ctx *gin.Context) {
	_, err := IsAuthorized(ctx.Request)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Error: Access token not found or not valid",
		})
		return
	}
	session, err := db.GetMongoSession()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error while getting mongo session",
		})
		return
	}
	defer session.Close()
	c := session.DB("dataImpact").C("users")
	var user usermodels.User
	userID := ctx.Param("id")
	err = c.Find(bson.M{"_id": userID}).One(&user)
	if err != nil {
		ctx.JSON(http.StatusNoContent, gin.H{
			"message": "id does not exist",
		})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

//Update: modify a user with its id, if the data field changes the file must be modified
func Update(ctx *gin.Context) {
	_, err := IsAuthorized(ctx.Request)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Error: Access token not found or not valid",
		})
		return
	}

	userID := ctx.Param("id")
	session, err := db.GetMongoSession()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error while getting mongo session",
		})
		return
	}
	defer session.Close()
	c := session.DB("dataImpact").C("users")

	var user usermodels.User
	err = ctx.Bind(&user)
	if err != nil {
		fmt.Println("Error = ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error while decoding",
		})
		return
	}
	err = c.Update(bson.M{"_id": userID}, user)
	if err != nil {
		ctx.JSON(http.StatusNoContent, gin.H{
			"message": "userID does not exist",
		})
		return
	}
	err = c.Find(bson.M{"_id": userID}).One(&user)
	if err != nil {
		ctx.JSON(http.StatusNoContent, gin.H{
			"message": "userID does not exist",
		})
		return
	}
	ctx.JSON(http.StatusOK, user)
}
