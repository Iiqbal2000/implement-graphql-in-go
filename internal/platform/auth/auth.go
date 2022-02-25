package platform

import (
	"time"

	"github.com/Iiqbal2000/mygopher/internal/users"
)

const TOKEN_ALIVE = time.Hour

type Auth struct {
	UserSvc users.UserService
}

func (a Auth) Authenticate(username, password string) (string, error) {
	userId, err := a.UserSvc.Compare(username, password)
	if err != nil {
		return "", err
	}

	token, err := CreateToken(userId, TOKEN_ALIVE)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a Auth) Authorize(tokenIn string) (string, error) {
	payload, err := VerifyToken(tokenIn)
	if err != nil {
		return "", err
	}

	return payload.UserID, nil
}
