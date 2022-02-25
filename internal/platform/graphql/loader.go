package platform

import (
	"context"

	"github.com/Iiqbal2000/mygopher"
	"github.com/Iiqbal2000/mygopher/internal/users"
	"github.com/graph-gophers/dataloader/v6"
)

type key string

var userLoaderKey key = "userLoader"

type Loaders struct {
	Collection map[key]dataloader.BatchFunc
}

func InitLoaders(userSvc users.UserService) Loaders {
	userL := userLoader{
		userSvc: userSvc,
	}

	return Loaders{
		Collection: map[key]dataloader.BatchFunc{
			userLoaderKey: userL.batchFunc,
		},
	}
}

// Attach all of batchF to context per-request
func (l Loaders) Attach(ctx context.Context) context.Context {
	for k, batchF := range l.Collection {
		ctx = context.WithValue(ctx, k, dataloader.NewBatchedLoader(batchF))
	}
	return ctx
}

func ExtractLoader(ctx context.Context, k key) (*dataloader.Loader, error) {
	result, ok := ctx.Value(k).(*dataloader.Loader)
	if !ok {
		return nil, mygopher.Error{
			Code: "500",
			Message: "internal server error",
		}
	}
	return result, nil
}