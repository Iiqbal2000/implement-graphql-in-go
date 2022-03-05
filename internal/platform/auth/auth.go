package platform

import (
	"time"

	"github.com/Iiqbal2000/mygopher/internal/users"
)

const TOKEN_ALIVE = time.Hour

type Auth struct {
	UserSvc users.UserService
}

func (a Auth) authenticate(username, password string) (string, error) {
	userId, err := a.UserSvc.Compare(username, password)
	if err != nil {
		return "", err
	}

	token, err := createToken(userId, TOKEN_ALIVE)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a Auth) authorize(tokenIn string) (string, error) {
	payload, err := verifyToken(tokenIn)
	if err != nil {
		return "", err
	}

	return payload.UserID, nil
}
