package main

import (
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	idCookieName    = "id"
	tokenCookieName = "Token"
	expDay          = 60 * 24
)

//Возвращает токен, считанный из куки
func (app *application) getTokenCookies(r *http.Request) *Token {
	cookieId, err := r.Cookie(idCookieName)
	if err != nil {
		return nil
	}
	cookieToken, err := r.Cookie(tokenCookieName)
	if err != nil {
		return nil
	}
	if cookieId.Value == "" || cookieToken.Value == "" {
		return nil
	}
	return &Token{
		IdUser: cookieId.Value,
		Token:  cookieToken.Value,
	}
}

//Возвращает новый объект куки
func newCookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:    name,
		Value:   value,
		Path:    "/",
		Expires: time.Now().Add(expDay * time.Hour),
	}
}

//Функция авторизации пользователя
//Ищет совпадения в базе пользователей
//Выдает новый токен доступа
//при успехе нет ошибки
func (app *application) auth(w http.ResponseWriter, email, password string) error {
	usr, err := getUserByEmail(email)
	if err != nil {
		return err
	}
	if usr == nil {
		return errors.New("user not found")
	}

	err = usr.comparePassword(password)
	if err != nil {
		return err
	}

	genToken, err := app.generateToken(usr.Id.Hex())
	if err != nil {
		return err
	}

	tkn := Token{
		IdUser: usr.Id.Hex(),
		Token:  genToken,
	}

	app.saveToken(w, *usr, tkn)
	return nil
}

//Проверка токена доступа, возвращает токен с данными и текущего пользователя при успехе
func (app *application) checkAuth(r *http.Request) (*Token, *User) {
	tkn := app.getTokenCookies(r)
	if tkn == nil {
		return nil, nil
	}

	//Декодируем токен из куки
	tDecode, err := base64.StdEncoding.DecodeString(tkn.Token)
	if err != nil {
		return nil, nil
	}

	tkn.Token = string(tDecode)
	usr := app.tokens.getUserByToken(*tkn)
	if usr == nil {
		return nil, nil
	}

	return tkn, usr
}

//Генерирует новый токен на основе слова
func (app *application) generateToken(word string) (string, error) {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	n := 20
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	bcryptB, err := bcrypt.GenerateFromPassword(b, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return word + strconv.FormatInt(time.Now().Unix(), 10) + string(bcryptB), nil
}
