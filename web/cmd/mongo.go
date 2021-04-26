package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	mongoUrl = "mongodb://localhost:27017"
	database = "web"
	usersCol = "users"
)

//Получение пользователя по адресу почты
func (app application) getUserByEmail(email string) (*User, error) {
	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	collection := session.DB(database).C(usersCol)
	var usr User
	err = collection.Find(bson.M{"email": email}).One(&usr)
	if err != nil {
		if err.Error() == "not found" {
			return nil, nil
		}
		return nil, err
	}
	return &usr, nil
}

//Возвращает список всех пользователей
func (app application) getAllUsers() ([]User, error) {
	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	collection := session.DB(database).C(usersCol)
	var users []User
	err = collection.Find(bson.M{}).All(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}
