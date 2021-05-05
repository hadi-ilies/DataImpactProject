package usercontrollers

import (
	"DataImpactProject/server/app/db"
	"DataImpactProject/server/app/models/usermodels"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

//Create: The JSON input file must be serialized and then saved to a MongoDB database. The JSON file will contain an array of users. You will then have to serialize the data concurrently and then insert it into the database without processing the already inserted entries again.The password must be encrypted with bcrypt and the hash of it inserted into the database.In addition to inserting this data into the database, you will have to generate one file per entry with the entry id as the filename, this file should contain only the Data field
func Create(ctx *gin.Context) {
	//read json request
	jsons, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while Reading request body",
		})
		return
	}
	//Note: I would use ctx.Bind this if the subject did not asked for desarialization concurrency
	//parse json array partially
	var serializedUser []interface{}
	err = json.Unmarshal(jsons, &serializedUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while Unmarshalling",
		})
		return
	}
	//init waitgroup
	wait := new(sync.WaitGroup)
	users := make([]usermodels.User, len(serializedUser))

	for i, s := range serializedUser {
		wait.Add(1)
		//goroutine for concurent Deserialization
		go func(j int, s interface{}) {
			defer wait.Done()
			deserializedObject, err := doDeserialization(s)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "Error while decoding",
				})
				return
			}
			users[j] = deserializedObject
		}(i, s)
	}
	wait.Wait()

	//save in db
	session, err := db.GetMongoSession()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while getting mongo session",
		})
		return
	}
	defer session.Close()
	c := session.DB("dataImpact").C("users")
	for _, user := range users {
		//check user email already exist
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
		//generate file with the user ID as filename and the data field as content
		err = createFile(user.Id, []byte(user.Data))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error while creating data file",
			})
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "users saved",
	})
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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while decoding",
		})
		return
	}
	session, err := db.GetMongoSession()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
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
	//check token
	_, err := IsAuthorized(ctx.Request)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Error: Access token not found or not valid",
		})
		return
	}

	session, err := db.GetMongoSession()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while getting mongo session",
		})
		return
	}
	defer session.Close()
	c := session.DB("dataImpact").C("users")
	userID := ctx.Param("id")
	err = c.Remove(bson.M{"_id": userID})
	if err != nil {
		ctx.JSON(http.StatusNoContent, gin.H{
			"message": "User does not exist",
		})
		return
	}
	err = deleteFile(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error: while deleting the file",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

//GetAllUsers: retrieve a list of users
func GetAllUsers(ctx *gin.Context) {
	//check token
	_, err := IsAuthorized(ctx.Request)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Error: Access token not found or not valid",
		})
		return
	}

	session, err := db.GetMongoSession()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
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
	//check token
	_, err := IsAuthorized(ctx.Request)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Error: Access token not found or not valid",
		})
		return
	}
	session, err := db.GetMongoSession()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
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
	//check token
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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while getting mongo session",
		})
		return
	}
	defer session.Close()
	c := session.DB("dataImpact").C("users")

	//get current user data
	var currentUser usermodels.User
	err = c.Find(bson.M{"_id": userID}).One(&currentUser)
	if err != nil {
		ctx.JSON(http.StatusNoContent, gin.H{
			"message": "userID does not exist",
		})
		return
	}

	//unmarshal data that will be updated
	var updatedUser usermodels.User
	err = ctx.Bind(&updatedUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error while decoding",
		})
		return
	}
	err = c.Update(bson.M{"_id": userID}, updatedUser)
	if err != nil {
		ctx.JSON(http.StatusNoContent, gin.H{
			"message": "userID does not exist",
		})
		return
	}
	//compare previous data field to the new one, If different we update the file otherwise we don't
	if currentUser.Data != updatedUser.Data {
		err = updateFile(userID, []byte(updatedUser.Data))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error While updating File",
			})
			return
		}
	}
	ctx.JSON(http.StatusOK, updatedUser)
}
