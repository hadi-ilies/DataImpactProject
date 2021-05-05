package usercontrollers

import (
	"DataImpactProject/server/app/models/usermodels"
	"io/ioutil"
	"os"
)

var directoryName string = "/tmp/"

//create File allows us to create a file as wanted inside the given subject
func createFile(filename string, data []byte) error {

	err := ioutil.WriteFile(directoryName+filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

//deleteFile allows us to delete a file as wanted inside the given subject
func deleteFile(filename string) error {
	err := os.Remove(directoryName + filename)
	if err != nil {
		return err
	}
	return nil
}

//updateFile: allows us to update a file as wanted inside the given subject (it's not mentionned that we need to keep the previous data therefore I just need to delete and recreate the file. and it's also bcs I'm super lazy)
func updateFile(filename string, data []byte) error {
	err := deleteFile(filename)
	if err != nil {
		return err
	}
	err = createFile(filename, data)
	if err != nil {
		return err
	}
	return nil
}

//userDesarialization: unserialize data with concurrency. we don't really handle errors here if an error happens we will send back an empty field
func userDesarialization(serializedUser map[string]interface{}) usermodels.User {
	var user usermodels.User

	if id, ok := serializedUser["id"].(string); ok {
		user.Id = id
	}
	if password, ok := serializedUser["password"].(string); ok {
		user.Password = password
	}
	if isActive, ok := serializedUser["isActive"].(bool); ok {
		user.IsActive = isActive
	}
	if balance, ok := serializedUser["balance"].(string); ok {
		user.Balance = balance
	}
	if age, ok := serializedUser["age"]; ok {
		user.Age = age
	}
	if name, ok := serializedUser["name"].(string); ok {
		user.Name = name
	}
	if gender, ok := serializedUser["gender"].(string); ok {
		user.Gender = gender
	}
	if company, ok := serializedUser["company"].(string); ok {
		user.Company = company
	}
	if email, ok := serializedUser["email"].(string); ok {
		user.Email = email
	}
	if phone, ok := serializedUser["phone"].(string); ok {
		user.Phone = phone
	}
	if address, ok := serializedUser["address"].(string); ok {
		user.Address = address
	}
	if about, ok := serializedUser["about"].(string); ok {
		user.About = about
	}
	if registred, ok := serializedUser["registred"].(string); ok {
		user.Registred = registred
	}
	if latitude, ok := serializedUser["latitude"].(float64); ok {
		user.Latitude = latitude
	}
	if longitude, ok := serializedUser["longitude"].(float64); ok {
		user.Longitude = longitude
	}
	if data, ok := serializedUser["data"].(string); ok {
		user.Data = data
	}

	//unsserialize tags
	for _, tag := range serializedUser["tags"].([]interface{}) {
		user.Tags = append(user.Tags, tag.(string))
	}
	//unserialize friends
	for _, serializedFriend := range serializedUser["friends"].([]interface{}) {
		var friend usermodels.Friend

		if name, ok := serializedFriend.(map[string]interface{})["name"].(string); ok {
			friend.Name = name
		}
		if id, ok := serializedFriend.(map[string]interface{})["id"].(uint64); ok {
			friend.Id = id
		}
		user.Friends = append(user.Friends, friend)
	}
	return user
}

//TODO ask DataImpact why Should I do that (desarialization concurrency)
//doDeserialization unmarshal my users with concurency .Runs inside goroutine
func doDeserialization(serializedUser interface{}) (usermodels.User, error) {
	m := serializedUser.(map[string]interface{})
	user := userDesarialization(m)

	return user, nil
}
