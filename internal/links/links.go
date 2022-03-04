package links

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strconv"

	"github.com/Iiqbal2000/mygopher"
	"github.com/Iiqbal2000/mygopher/internal/users"
	"github.com/huandu/go-sqlbuilder"
)

var (
	errBadTitle           = errors.New("the title value is bad")
	errBadAddress         = errors.New("the address value is bad")
	errBadUserId          = errors.New("the user id value is bad")
	errLinkInternalServer = errors.New("internal server error")
)

type LinkService struct {
	Db      *sql.DB
	Log     *log.Logger
	UserSvc users.UserService
}

func (l LinkService) CreateLink(ctx context.Context, linkIn mygopher.LinkIn) (mygopher.LinkOut, error) {
	isAuthorized := ctx.Value("is_authorized").(bool)
	if !isAuthorized {
		return mygopher.LinkOut{}, mygopher.Error{
			Code:    "401",
			Message: "Unauthorized",
		}
	}

	userId := ctx.Value("user_id").(string)
	linkIn.UserID = userId

	switch {
	case linkIn.Title == "":
		return mygopher.LinkOut{}, mygopher.Error{
			Code:    "400",
			Message: errBadTitle.Error(),
		}
	case linkIn.Address == "":
		return mygopher.LinkOut{}, mygopher.Error{
			Code:    "400",
			Message: errBadAddress.Error(),
		}
	case linkIn.UserID == "":
		return mygopher.LinkOut{}, mygopher.Error{
			Code:    "400",
			Message: errBadUserId.Error(),
		}
	}

	link := mygopher.Link{
		Title:   linkIn.Title,
		Address: linkIn.Address,
		UserID:  linkIn.UserID,
	}

	query, args := sqlbuilder.InsertInto("links").
		Cols("title", "address", "userId").
		Values(link.Title, link.Address, link.UserID).
		Build()

	result, err := l.Db.Exec(query, args...)
	if err != nil {
		l.Log.Println("failure when executing link query")
		return mygopher.LinkOut{}, mygopher.Error{
			Code:    "500",
			Message: "internal server error",
		}
	}

	linkIdRaw, err := result.LastInsertId()
	if err != nil {
		l.Log.Println("failure when get link id: ", err.Error())
		return mygopher.LinkOut{}, mygopher.Error{
			Code:    "500",
			Message: "internal server error",
		}
	}

	linkId := strconv.FormatInt(linkIdRaw, 10)

	return mygopher.LinkOut{
		ID:      linkId,
		Title:   link.Title,
		Address: link.Address,
		UserID:  link.UserID,
	}, nil
}

func (l LinkService) GetAll(ctx context.Context) ([]mygopher.LinkOut, error) {
	isAuthorized := ctx.Value("is_authorized").(bool)
	if !isAuthorized {
		return nil, mygopher.Error{
			Code:    "401",
			Message: "Unauthorized",
		}
	}

	query, args := sqlbuilder.Select("rowid", "title", "address", "userId").
		From("links").
		Build()

	rows, err := l.Db.Query(query, args...)
	if err != nil {
		l.Log.Println("failure when get links: ", err.Error())
		return nil, mygopher.Error{
			Code:    "500",
			Message: errLinkInternalServer.Error(),
		}
	}

	defer rows.Close()

	linksOut := make([]mygopher.LinkOut, 0)

	for rows.Next() {
		link := mygopher.LinkOut{}
		var linkId int64
		var userId int64

		err := rows.Scan(&linkId, &link.Title, &link.Address, &userId)
		if err != nil {
			l.Log.Println("failure when in iteration links: ", err.Error())
			return nil, mygopher.Error{
				Code:    "500",
				Message: errLinkInternalServer.Error(),
			}
		}

		link.ID = strconv.FormatInt(linkId, 10)
		link.UserID = strconv.FormatInt(userId, 10)

		linksOut = append(linksOut, link)
	}

	return linksOut, nil
}
