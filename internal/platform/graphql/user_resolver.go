package platform

import (
	"github.com/Iiqbal2000/mygopher"
	"github.com/graph-gophers/graphql-go"
)

type UserResolver struct {
	User mygopher.UserOut
}

func (u *UserResolver) Id() graphql.ID {
	return graphql.ID(u.User.ID)
}

func (u *UserResolver) Name() string {
	return u.User.Username
}
