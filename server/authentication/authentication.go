package authentication

import "strings"

type UserManager struct {
	users map[string]string
}

func (u *UserManager) ValidateUser(username string, password string) (string, bool) {
	for k, v := range u.users {
		if strings.EqualFold(k, username) && password == v {
			return k, true
		}
	}
	return "", false
}

func (u *UserManager) RegisterUser(username string, password string) {
	u.users[username] = password
}

/*
@returns if a username was unregistered
 */
func (u *UserManager) UnregisterUser(username string) bool {
	_, has := u.users[username]
	if has {
		delete(u.users, username)
	}
	return has
}

func NewUserManager() UserManager {
	return UserManager{make(map[string]string)}
}
