package main

type UserAndTokens struct {
	User   User
	Tokens []Token
}

//Карта, где хранятся все токены и  авторизированные пользователи
//Ключ - это id пользователя, значение - это пользователь и его токены
type MapTokens map[string]*UserAndTokens

func newMapTokens() *MapTokens {
	result := make(MapTokens)
	return &result
}

//Обновление, передается пользователь, если он есть в карте, то обновляются его данные
func (m *MapTokens) updateUser(usr User) {
	id := usr.Id.Hex()
	//Если его нет в карте, то ничего не делается
	if (*m)[id] == nil {
		return
	}
	(*m)[id].User = usr
}

//Добавление, передается пользователь и токен
func (m *MapTokens) add(usr User, tkn Token) {
	id := usr.Id.Hex()
	//Если этого пользователя нет в карте, создается новая запись с ним
	if (*m)[id] == nil {
		(*m)[id] = &UserAndTokens{
			User:   usr,
			Tokens: []Token{tkn},
		}
		//Если же он уже есть, то к списку его токенов добавляется новый
	} else {
		(*m)[id].User = usr
		(*m)[id].Tokens = append((*m)[id].Tokens, tkn)
	}
}

//Удаляются все токены и сам пользователь из карты
func (m *MapTokens) clearById(id string) {
	delete(*m, id)
}

//Удаляет токен пользователя из карты
func (m *MapTokens) deleteByToken(tkn Token) {
	id := tkn.IdUser

	//Если записи нет, то ничего не делаем
	if (*m)[id] == nil {
		return
	}

	//Пересобираем токены без учета удаляемого
	var newSlice []Token
	for _, el := range (*m)[id].Tokens {
		if el.Token != tkn.Token {
			newSlice = append(newSlice, el)
		}
	}
	(*m)[id].Tokens = newSlice

	//Если не осталось токенов, то удаляем запись в карте
	if len(newSlice) == 0 {
		delete(*m, id)
	}
}

//Получаем хозяина токена, если он есть
func (m MapTokens) getUserByToken(tkn Token) *User {
	id := tkn.IdUser

	//Если нет записи этого пользователя
	if m[id] == nil {
		return nil
	}

	//Перебираем все токены пользователя на предмет совпадения, чтобы вернуть искомый
	for _, el := range m[id].Tokens {
		if el.Token == tkn.Token {
			return &m[id].User
		}
	}
	return nil
}
