package main

//Карта, где хранятся все токены и  авторизированные пользователи
type mapTokens map[string] *user

func newMapTokens() *mapTokens {
	result := make(mapTokens)
	return &result
}

func (m *mapTokens) save(u *user, token string) {
	(*m)[token] = u
}

func (m mapTokens) getUser(token string) *user {
	return m[token]
}

func (m *mapTokens) remove(token string) {
	delete(*m, token)
}