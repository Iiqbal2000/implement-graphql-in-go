package users

import (
	"crypto/subtle"
	"database/sql"
	"errors"
	"log"
	"strconv"

	"github.com/Iiqbal2000/mygopher"
	"github.com/huandu/go-sqlbuilder"
)

var (
	errBadUsername  = errors.New("the username value is bad")
	errBadPassword  = errors.New("the password value is bad")
	errUserNotFound = errors.New("the user was not found")
	errWrongPass    = errors.New("the password does not match")
)

type Input struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Output struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type UserService struct {
	Db  *sql.DB
	Log *log.Logger
}

func (u UserService) Add(input Input) (Output, error) {
	if input.Username == "" {
		return Output{}, mygopher.Error{
			Code:    "400",
			Message: errBadUsername.Error(),
		}
	} else if input.Password == "" {
		return Output{}, mygopher.Error{
			Code:    "400",
			Message: errBadPassword.Error(),
		}
	}

	user := mygopher.User{
		Username: input.Username,
		Password: input.Password,
	}

	query, args := sqlbuilder.InsertInto("users").
		Cols("username", "password").
		Values(user.Username, user.Password).
		Build()

	result, err := u.Db.Exec(query, args...)
	if err != nil {
		u.Log.Println("failure when inserting new user")
		return Output{}, mygopher.Error{
			Code:    "500",
			Message: "internal server error",
		}
	}

	userIdRaw, err := result.LastInsertId()
	if err != nil {
		u.Log.Println("failure when getting user id")
		return Output{}, mygopher.Error{
			Code:    "500",
			Message: "internal server error",
		}
	}

	userId := strconv.FormatInt(userIdRaw, 10)

	return Output{
		ID:       userId,
		Username: user.Username,
	}, nil
}

func (u UserService) GetByIds(ids []string) ([]Output, error) {
	sb := sqlbuilder.NewSelectBuilder()

	// []int{1,2} => "?,?"
	listMark := sqlbuilder.List(ids)

	query, args := sb.Select("rowid", "username").
		From("users").
		Where(sb.In("rowid", listMark)).
		Build()

	u.Log.Println(query)

	rows, err := u.Db.Query(query, args...)
	if err != nil {
		u.Log.Println("failure when querying users: ", err.Error())
		return nil, mygopher.Error{
			Code:    "500",
			Message: "internal server error",
		}
	}

	defer rows.Close()

	users := make([]Output, len(ids))

	for rows.Next() {
		user := Output{}
		var userId int64

		err := rows.Scan(&userId, &user.Username)
		if err != nil {
			u.Log.Println("failure in iteraiton of users: ", err.Error())
			return nil, mygopher.Error{
				Code:    "500",
				Message: "internal server error",
			}
		}

		user.ID = strconv.FormatInt(userId, 10)
		users = append(users, user)
	}

	return users, nil
}

func (u UserService) GetById(id string) (Output, error) {
	user := Output{}
	var userId int64

	sb := sqlbuilder.NewSelectBuilder()
	query, args := sb.Select("rowid, username").
		From("users").
		Where(sb.Equal("rowid", id)).
		Build()

	err := u.Db.QueryRow(query, args...).Scan(&userId, &user.Username)
	if err != nil {
		u.Log.Println("failure when getting a user: ", err.Error())
		return Output{}, mygopher.Error{
			Code:    "500",
			Message: "internal server error",
		}
	}

	user.ID = strconv.FormatInt(userId, 10)

	return user, nil
}

func (u UserService) GetByUsername(username string) (mygopher.User, error) {
	user := mygopher.User{}
	var userId int64

	sb := sqlbuilder.NewSelectBuilder()
	query, arg := sb.Select("rowid", "username", "password").
		From("users").
		Where(sb.Equal("username", username)).
		Build()

	err := u.Db.QueryRow(query, arg...).Scan(&userId, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return mygopher.User{}, mygopher.Error{
				Code:    "404",
				Message: errUserNotFound.Error(),
			}
		}

		u.Log.Println("failure when getting a user by username: ", err.Error())
		return mygopher.User{}, mygopher.Error{
			Code:    "500",
			Message: "internal server error",
		}
	}

	user.ID = strconv.FormatInt(userId, 10)

	return user, nil
}

func (u UserService) Compare(usernameIn, passwordIn string) (string, error) {
	// check whether the username exists
	user, err := u.GetByUsername(usernameIn)
	if err != nil {
		return "", err
	}

	// comparing the password
	passwordMatch := (subtle.ConstantTimeCompare([]byte(passwordIn), []byte(user.Password)) == 1)
	if !passwordMatch {
		return "", mygopher.Error{
			Code:    "400",
			Message: errWrongPass.Error(),
		}
	}

	return user.ID, nil
}
