package platform

import (
	"time"

	"github.com/Iiqbal2000/mygopher/internal/users"
)

const TOKEN_ALIVE = time.Hour

type Service struct {
	UserSvc users.Service
}

func (s Service) authenticate(username, password string) (string, error) {
	userId, err := s.UserSvc.Compare(username, password)
	if err != nil {
		return "", err
	}

	token, err := createToken(userId, TOKEN_ALIVE)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s Service) authorize(tokenIn string) (string, error) {
	payload, err := verifyToken(tokenIn)
	if err != nil {
		return "", err
	}

	return payload.UserID, nil
}
