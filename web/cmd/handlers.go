package main

import (
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
)

//Главная страница
func (app *application) indexPageGET(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	ts, err := template.ParseFiles(
		"./ui/views/page.index.tmpl",
		"./ui/views/header.main.tmpl",
		"./ui/views/footer.main.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}

	tkn, usr := app.checkAuth(r)
	if tkn == nil || usr == nil {
		err = ts.Execute(w, struct {
			User *User
		}{
			User: nil,
		})
	} else {
		err = ts.Execute(w, struct {
			User *User
		}{
			User: usr,
		})
	}
	if err != nil {
		app.serverError(w, err)
	}
}

//Страница отображения всех пользователей
func (app *application) usersPageGET(w http.ResponseWriter, r *http.Request) {
	tkn, usr := app.checkAuth(r)
	if tkn == nil || usr == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	ts, err := template.ParseFiles(
		"./ui/views/page.users.tmpl",
		"./ui/views/header.main.tmpl",
		"./ui/views/footer.main.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}

	users, err := app.getAllUsers()
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, struct {
		User  *User
		Users []User
	}{
		User:  usr,
		Users: users,
	})
	if err != nil {
		app.serverError(w, err)
	}
}

//Отображение страницы регистрации
func (app *application) signUpPageGET(w http.ResponseWriter, r *http.Request) {
	tkn, usr := app.checkAuth(r)
	if tkn != nil || usr != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	ts, err := template.ParseFiles(
		"./ui/views/page.signUp.tmpl",
		"./ui/views/header.main.tmpl",
		"./ui/views/footer.main.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

//Обработка POST-запроса страницы регистрации
func (app *application) signUpPagePOST(w http.ResponseWriter, r *http.Request) {
	usr := User{
		Id:       bson.NewObjectId(),
		Email:    r.FormValue("email"),
		Name:     r.FormValue("name"),
		Surname:  r.FormValue("surname"),
		Password: r.FormValue("password"),
	}
	repPassword := r.FormValue("repPassword")

	valid, err := app.valid(usr, repPassword)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if !valid {
		http.Redirect(w, r, "/signUp/", http.StatusSeeOther)
		return
	}

	err = app.saveUser(usr)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.infoLog.Println("Новый пользователь:", usr.Email)
	http.Redirect(w, r, "/signIn/", http.StatusSeeOther)
}

//Отображение страницы авторизации
func (app *application) signInPageGET(w http.ResponseWriter, r *http.Request) {
	tkn, usr := app.checkAuth(r)
	if tkn != nil || usr != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	ts, err := template.ParseFiles(
		"./ui/views/page.signIn.tmpl",
		"./ui/views/header.main.tmpl",
		"./ui/views/footer.main.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

//Обработка POST-запроса страницы авторизации
func (app *application) signInPagePOST(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		http.Redirect(w, r, "/signIn/", http.StatusSeeOther)
		return
	}

	//auth user
	usr, err := app.getUserByEmail(email)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if usr == nil || usr.comparePassword(password) != nil {
		http.Redirect(w, r, "/signIn/", http.StatusSeeOther)
		return
	}
	genToken, err := app.generateToken(usr.Id.Hex())
	if err != nil {
		app.serverError(w, err)
		return
	}

	tkn := Token{
		IdUser: usr.Id.Hex(),
		Token:  genToken,
	}
	app.saveToken(w, *usr, tkn)
	app.infoLog.Println("Пользователь вошел:", email, "\tid:", usr.Id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//Выход из учетной записи
func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	tkn, usr := app.checkAuth(r)
	if tkn == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	app.deleteToken(w, *usr, *tkn)
	app.infoLog.Println("Пользователь вышел:", usr.Email, "\tid:", usr.Id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//Страница изменения пользователя
func (app *application) changeUserGET(w http.ResponseWriter, r *http.Request) {
	tkn, usr := app.checkAuth(r)
	if tkn == nil || usr == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	ts, err := template.ParseFiles(
		"./ui/views/page.changeUser.tmpl",
		"./ui/views/header.main.tmpl",
		"./ui/views/footer.main.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, struct {
		User *User
	}{
		User: usr,
	})
	if err != nil {
		app.serverError(w, err)
	}
}

//Обработка запроса на смену данных пользователя
func (app *application) changeUserPOST(w http.ResponseWriter, r *http.Request) {
	tkn, usr := app.checkAuth(r)
	if tkn == nil || usr == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	newU := User{
		Id:       usr.Id,
		Email:    r.FormValue("email"),
		Name:     r.FormValue("name"),
		Surname:  r.FormValue("surname"),
		Password: usr.Password,
	}
	valid, err := app.valid(newU, newU.Password)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if !valid {
		http.Redirect(w, r, "/changeUser/", http.StatusSeeOther)
		return
	}

	err = app.updateUser(newU)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.tokens.updateUser(newU)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//Отображение страницы смены пароля
func (app *application) changePasswordGET(w http.ResponseWriter, r *http.Request) {
	tkn, usr := app.checkAuth(r)
	if tkn == nil || usr == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	ts, err := template.ParseFiles(
		"./ui/views/page.changePassword.tmpl",
		"./ui/views/header.main.tmpl",
		"./ui/views/footer.main.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, struct {
		User *User
	}{
		User: usr,
	})
	if err != nil {
		app.serverError(w, err)
	}
}

//Обработка запроса на обновление пароля пользователя
func (app *application) changePasswordPOST(w http.ResponseWriter, r *http.Request) {
	tkn, usr := app.checkAuth(r)
	if tkn == nil || usr == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	password := r.FormValue("password")
	newPassword := r.FormValue("newPassword")
	repNewPassword := r.FormValue("repNewPassword")
	if password == "" || newPassword == "" || repNewPassword == "" ||
		newPassword != repNewPassword {
		http.Redirect(w, r, "/changePassword/", http.StatusSeeOther)
		return
	}

	err := usr.comparePassword(password)
	if err != nil {
		http.Redirect(w, r, "/changePassword/", http.StatusSeeOther)
		return
	}

	err = app.updateUserPassword(*usr, newPassword)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.tokens.updateUser(*usr)
	http.Redirect(w, r, "/logout/", http.StatusSeeOther)
}

//Отображение страницы удаления пользователя
func (app *application) deleteUserGET(w http.ResponseWriter, r *http.Request) {
	tkn, usr := app.checkAuth(r)
	if tkn == nil || usr == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	ts, err := template.ParseFiles(
		"./ui/views/page.deleteUser.tmpl",
		"./ui/views/header.main.tmpl",
		"./ui/views/footer.main.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, struct {
		User *User
	}{
		User: usr,
	})
	if err != nil {
		app.serverError(w, err)
	}
}

//Обработка запроса на удаление пользователя
func (app *application) deleteUserPOST(w http.ResponseWriter, r *http.Request) {
	tkn, usr := app.checkAuth(r)
	if tkn == nil || usr == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	if password == "" || email == "" {
		http.Redirect(w, r, "/deleteUser/", http.StatusSeeOther)
		return
	}

	err := usr.comparePassword(password)
	if err != nil {
		http.Redirect(w, r, "/deleteUser/", http.StatusSeeOther)
		return
	}

	if email != usr.Email {
		http.Redirect(w, r, "/deleteUser/", http.StatusSeeOther)
		return
	}

	err = app.deleteUser(*usr)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.tokens.clearById(tkn.IdUser)
	http.Redirect(w, r, "/logout/", http.StatusSeeOther)
}
