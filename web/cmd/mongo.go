package main

import (
	"golang.org/x/crypto/bcrypt"
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

//Сохранение пользователя в базе
func (app application) saveUser(usr User) error {
	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		return err
	}
	defer session.Close()

	collection := session.DB(database).C(usersCol)
	bcryptPassw, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	usr.Password = string(bcryptPassw)
	usr.Id = bson.NewObjectId()
	err = collection.Insert(usr)
	if err != nil {
		return err
	}
	return nil
}

//Обновление данных пользователя
func (app application) updateUser(usr User) error {
	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		return err
	}
	defer session.Close()

	collection := session.DB(database).C(usersCol)
	err = collection.Update(bson.M{"_id": usr.Id}, usr)
	if err != nil {
		return err
	}
	return nil
}

//Обновление пароля пользователя
func (app application) updateUserPassword(usr User, password string) error {
	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		return err
	}
	defer session.Close()

	collection := session.DB(database).C(usersCol)
	bcryptPassw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	usr.Password = string(bcryptPassw)
	err = collection.Update(bson.M{"_id": usr.Id}, usr)
	if err != nil {
		return err
	}
	return nil
}

//Удаление пользователя из базы
func (app application) deleteUser(usr User) error {
	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		return err
	}
	defer session.Close()

	collection := session.DB(database).C(usersCol)
	err = collection.Remove(bson.M{"_id": usr.Id})
	if err != nil {
		return err
	}
	return nil
}
