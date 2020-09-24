package auth_service

import "material/models"

type Auth struct {
	Username string
	Password string
}

func (a *Auth) Check() (int64, error) {
	return models.CheckAuth(a.Username, a.Password)
}
