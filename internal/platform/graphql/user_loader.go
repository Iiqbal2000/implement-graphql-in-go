package platform

import (
	"context"

	"github.com/Iiqbal2000/mygopher"
	"github.com/Iiqbal2000/mygopher/internal/users"
	"github.com/graph-gophers/dataloader/v6"
)

type userLoader struct {
	userSvc users.UserService
}

func (u userLoader) batchFunc(_ context.Context, keys dataloader.Keys) []*dataloader.Result {
	var results []*dataloader.Result
	batch := make(map[string]mygopher.UserOut)

	users, err := u.userSvc.GetByIds(keys.Keys())
	for _, user := range users {
		batch[user.ID] = user
	}

	for _, key := range keys {
		results = append(results, &dataloader.Result{
			Data:  batch[key.String()],
			Error: err,
		})
	}

	return results
}
