package platform

import (
	"context"

	"github.com/Iiqbal2000/mygopher"
	"github.com/graph-gophers/dataloader/v6"
	"github.com/graph-gophers/graphql-go"
)

type LinkResolver struct {
	LinkOut    mygopher.LinkOut
	UserOut mygopher.UserOut
	loader  *dataloader.Loader
}

func (l *LinkResolver) ID() graphql.ID {
	return graphql.ID(l.LinkOut.ID)
}

func (l *LinkResolver) Title() string {
	return l.LinkOut.Title
}

func (l *LinkResolver) Address() string {
	return l.LinkOut.Address
}

func (l *LinkResolver) User() (*UserResolver, error) {
	// if l.loader == nil {
	// 	return &UserResolver{User: l.UserOut}, nil
	// }

	result, err := l.loader.Load(context.TODO(), dataloader.StringKey(l.LinkOut.UserID))()
	if err != nil {
		return &UserResolver{}, err
	}

	out := result.(mygopher.UserOut)
	return &UserResolver{User: out}, nil
}
