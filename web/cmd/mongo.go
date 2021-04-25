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
func getUserByEmail(email string) (*user, error) {
	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		return nil, err
	}
	defer session.Close()
	collection := session.DB(database).C(usersCol)
	var u user
	err = collection.Find(bson.M{"email": email}).One(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func getUserById(id bson.ObjectId) (*user, error) {
	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		return nil, err
	}
	defer session.Close()
	collection := session.DB(database).C(usersCol)
	var u user
	err = collection.Find(bson.M{"_id": id}).One(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

//Возвращает список всех пользователей
func getAllUsers() ([]user, error) {
	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		return nil, err
	}
	defer session.Close()
	collection := session.DB(database).C(usersCol)
	var users []user
	err = collection.Find(bson.M{}).All(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}
