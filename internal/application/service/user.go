package service

import (
	"context"
	"log"
	"net/http"

	"github.com/fedo3nik/GamePortal_IdentityService/internal/application"
	"github.com/fedo3nik/GamePortal_IdentityService/internal/domain/entities"
	e "github.com/fedo3nik/GamePortal_IdentityService/internal/error"
	"github.com/fedo3nik/GamePortal_IdentityService/internal/infrastructure/database/mongodb"
	"github.com/gorilla/securecookie"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User interface {
	SignUp(ctx context.Context, nickname, email, password string) (*entities.User, error)
	SignIn(ctx context.Context, email, password string) (*entities.User, *http.Cookie, error)
	AddWarning(ctx context.Context, id string) (*entities.User, error)
	RemoveWarning(ctx context.Context, id string) (*entities.User, error)
	IsBanned(ctx context.Context, id string) (*entities.User, bool, error)
}

type UserService struct {
	client *mongo.Client
	db     string
}

func (u UserService) SignUp(ctx context.Context, nickname, email, password string) (*entities.User, error) {
	if !application.IsEmailValid(email) {
		return nil, errors.Wrap(e.ErrValidation, "invalid email")
	}

	if !application.IsNicknameValid(nickname) {
		return nil, errors.Wrap(e.ErrValidation, "invalid nickname")
	}

	if !application.IsPasswordValid(password) {
		return nil, errors.Wrap(e.ErrValidation, "invalid password")
	}

	var checkUser = true

	var usr entities.User

	usersCollection := mongodb.GetCollection(u.client, u.db)

	_, err := mongodb.GetDocumentByEmail(ctx, usersCollection, email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			checkUser = false
		} else {
			return nil, err
		}
	}

	if checkUser {
		return nil, errors.Wrap(e.ErrDB, "user with this email already exist")
	}

	usr.Nickname = nickname
	usr.Email = email
	usr.Password = password

	result, err := mongodb.Insert(ctx, usersCollection, &usr)
	if err != nil {
		log.Printf("Insert document into db err: %v", err)
		return nil, errors.Wrap(e.ErrDB, "insert document")
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		usr.ID = oid.Hex()
	} else {
		return nil, e.ErrDB
	}

	return &usr, nil
}

func (u UserService) SignIn(ctx context.Context, email, password string) (*entities.User, *http.Cookie, error) {
	if !application.IsEmailValid(email) {
		return nil, nil, errors.Wrap(e.ErrValidation, "invalid email")
	}

	if !application.IsPasswordValid(password) {
		return nil, nil, errors.Wrap(e.ErrValidation, "invalid password")
	}

	var usr entities.User

	usersCollection := mongodb.GetCollection(u.client, u.db)

	user, err := mongodb.GetDocumentByEmail(ctx, usersCollection, email)
	if err != nil {
		log.Printf("Get document from db err: %v", err)
		return nil, nil, err
	}

	if password == user.Password {
		usr.ID = user.ID.Hex()
		usr.Email = user.Email
		usr.Password = user.Password
		usr.Nickname = user.Nickname
		usr.WarningCount = user.WarningCount

		hashKey := []byte(email + password)
		s := securecookie.New(hashKey, nil)

		encoded, errCookie := s.Encode("user-cookie", usr.Nickname)
		if errCookie != nil {
			log.Printf("Secure cookie encoded error: %v", errCookie)
			return nil, nil, errCookie
		}

		cookie := &http.Cookie{
			Name:  "user-cookie",
			Value: encoded,
			Path:  "/user/sign-in",
		}

		return &usr, cookie, nil
	}

	return nil, nil, errors.Wrap(e.ErrSignIn, "wrong password")
}

func (u UserService) AddWarning(ctx context.Context, id string) (*entities.User, error) {
	var usr entities.User

	usersCollection := mongodb.GetCollection(u.client, u.db)

	user, err := mongodb.GetDocumentByID(ctx, usersCollection, id)
	if err != nil {
		log.Printf("Get document from db err: %v", err)
		return nil, e.ErrDB
	}

	updateResult, err := mongodb.UpdateWarningCountField(ctx, usersCollection, id, user.WarningCount+1)
	if err != nil {
		return nil, err
	}

	if updateResult.ModifiedCount == 0 {
		return nil, errors.Wrap(e.ErrDB, "update error")
	}

	usr.ID = user.ID.Hex()
	usr.Email = user.Email
	usr.Password = user.Password
	usr.Nickname = user.Nickname
	usr.WarningCount = user.WarningCount + 1

	return &usr, nil
}

func (u UserService) RemoveWarning(ctx context.Context, id string) (*entities.User, error) {
	var usr entities.User

	usersCollection := mongodb.GetCollection(u.client, u.db)

	user, err := mongodb.GetDocumentByID(ctx, usersCollection, id)
	if err != nil {
		log.Printf("Get document from db err: %v", err)
		return nil, err
	}

	if user.WarningCount == 0 {
		usr.WarningCount = user.WarningCount
		return &usr, nil
	}

	updateResult, err := mongodb.UpdateWarningCountField(ctx, usersCollection, id, user.WarningCount-1)
	if err != nil {
		return nil, err
	}

	if updateResult.ModifiedCount == 0 {
		return nil, errors.Wrap(e.ErrDB, "update error")
	}

	usr.ID = user.ID.Hex()
	usr.Email = user.Email
	usr.Password = user.Password
	usr.Nickname = user.Nickname
	usr.WarningCount = user.WarningCount - 1

	return &usr, nil
}

func (u UserService) IsBanned(ctx context.Context, id string) (*entities.User, bool, error) {
	var usr entities.User

	var borderWarningCount uint = 3

	usersCollection := mongodb.GetCollection(u.client, u.db)

	user, err := mongodb.GetDocumentByID(ctx, usersCollection, id)
	if err != nil {
		log.Printf("Get document from db err: %v", err)
		return nil, false, err
	}

	usr.ID = user.ID.Hex()
	usr.Email = user.Email
	usr.Password = user.Password
	usr.Nickname = user.Nickname
	usr.WarningCount = user.WarningCount

	if user.WarningCount >= borderWarningCount {
		return &usr, true, nil
	}

	return &usr, false, nil
}

func NewUserService(client *mongo.Client, db string) *UserService {
	return &UserService{client: client, db: db}
}
