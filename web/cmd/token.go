package main

import (
	"encoding/base64"
	"net/http"
)

type token struct {
	IdUser string
	Token  string
}

//Проверка токена на пустоту
func (t token) isEmpty() bool {
	if t.Token == "" || t.IdUser == "" {
		return true
	}
	return false
}

//Сохраняет токен
func (app *application) saveToken(w http.ResponseWriter, t token, u *user) {
	http.SetCookie(w,
		newCookie(idCookieName, t.IdUser))
	//base64 token save in cookie
	base64Tkn := base64.StdEncoding.EncodeToString([]byte(t.Token))
	http.SetCookie(w,
		newCookie(tokenCookieName, base64Tkn))
	app.tokens.save(u, t.Token)
}

//Удаляет токен
func (app *application) deleteToken(w http.ResponseWriter, t string) error {
	http.SetCookie(w,
		newCookie(idCookieName, ""))
	http.SetCookie(w,
		newCookie(tokenCookieName, ""))
	//Декодируем токен из куки
	tDecode, err := base64.StdEncoding.DecodeString(t)
	if err != nil {
		return err
	}
	app.tokens.remove(string(tDecode))
	return nil
}
