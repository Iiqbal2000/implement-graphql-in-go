package platform

import (
	"github.com/Iiqbal2000/mygopher/internal/users"
	"github.com/graph-gophers/graphql-go"
)

type UserResolver struct {
	User users.Output
}

func (u *UserResolver) Id() graphql.ID {
	return graphql.ID(u.User.ID)
}

func (u *UserResolver) Name() string {
	return u.User.Username
}
