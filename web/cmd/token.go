package main

import (
	"encoding/base64"
	"net/http"
)

type Token struct {
	IdUser string
	Token  string
}

//Сохраняет токен
func (app *application) saveToken(w http.ResponseWriter, u User, tkn Token) {
	http.SetCookie(w,
		newCookie(idCookieName, tkn.IdUser))

	//base64 Token save in cookie
	base64Tkn := base64.StdEncoding.EncodeToString([]byte(tkn.Token))
	http.SetCookie(w,
		newCookie(tokenCookieName, base64Tkn))

	app.tokens.add(u, tkn)
}

//Удаляет токен
func (app *application) deleteToken(w http.ResponseWriter, u User, tkn Token) {
	http.SetCookie(w,
		newCookie(idCookieName, ""))
	http.SetCookie(w,
		newCookie(tokenCookieName, ""))

	app.tokens.deleteByToken(tkn)
}
