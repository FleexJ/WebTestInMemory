package main

import (
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
	"regexp"
)

const regexEmail = `^\w+@\w+[.]\w+$`

type User struct {
	Id       bson.ObjectId `bson:"_id"`
	Email    string
	Name     string
	Surname  string
	Password string
}

//Валидация пользователя перед записью в базу
func (app application) valid(usr User, repPassword string) (bool, error) {
	matched, _ := regexp.MatchString(regexEmail, usr.Email)
	if !matched ||
		usr.Name == "" ||
		usr.Surname == "" ||
		usr.Password == "" ||
		usr.Password != repPassword {
		return false, nil
	}
	uG, err := app.getUserByEmail(usr.Email)
	if err != nil {
		return false, err
	}
	if uG != nil && usr.Id != uG.Id {
		return false, nil
	}
	return true, nil
}

//Сравнение пароля пользователя
func (usr User) comparePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
