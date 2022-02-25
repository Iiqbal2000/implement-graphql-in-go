package platform

import (
	"context"

	"github.com/Iiqbal2000/mygopher"
	"github.com/Iiqbal2000/mygopher/internal/links"
	"github.com/Iiqbal2000/mygopher/internal/users"
)

// the root resolver
type Resolver struct {
	LinkSvc links.LinkService
	UserSvc users.UserService
}

func (r *Resolver) Links(ctx context.Context) (*[]*LinkResolver, error) {
	linkResolver := make([]*LinkResolver, 0)

	links, err := r.LinkSvc.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// extract a loader from context
	loader, err := ExtractLoader(ctx, userLoaderKey)
	if err != nil {
		return nil, err
	}

	for _, l := range links {
		linkResolver = append(linkResolver, &LinkResolver{
			LinkOut: l,
			loader:  loader,
		})
	}

	return &linkResolver, nil
}

func (r *Resolver) CreateLink(ctx context.Context, args struct{ Input mygopher.LinkIn }) (*LinkResolver, error) {
	linkOut, err := r.LinkSvc.CreateLink(ctx, args.Input)
	if err != nil {
		return nil, err
	}

	userOut, err := r.UserSvc.GetById(linkOut.UserID)
	if err != nil {
		return nil, err
	}

	return &LinkResolver{
		LinkOut: linkOut,
		UserOut: userOut,
	}, nil
}

func (r *Resolver) CreateUser(args struct{ Input mygopher.UserIn }) (*UserResolver, error) {
	result, err := r.UserSvc.Add(args.Input)
	if err != nil {
		return nil, err
	}

	return &UserResolver{result}, nil
}
