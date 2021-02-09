package main

import (
	"context"
	"testing"

	"github.com/fedo3nik/GamePortal_IdentityService/config"
	"github.com/fedo3nik/GamePortal_IdentityService/internal/application/service"
	e "github.com/fedo3nik/GamePortal_IdentityService/internal/error"
	"github.com/fedo3nik/GamePortal_IdentityService/internal/infrastructure/database/mongodb"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestIsBanned(t *testing.T) {
	c, err := config.NewConfig()
	assert.NoError(t, err)

	client, err := InitClient(context.Background(), c.ConnURI)
	assert.NoError(t, err)

	defer func() {
		err = client.Disconnect(context.Background())
		assert.NoError(t, err)
	}()

	collection := mongodb.GetCollection(client, c.DB)
	userService := service.NewUserService(client, c.DB)

	usr, err := userService.SignUp(context.Background(), "testCli", "testcli@test.com", "testCli123")
	assert.NoError(t, err)

	defer func() {
		_, err := mongodb.DeleteAll(context.Background(), collection)
		assert.NoError(t, err)
	}()

	tests := []struct {
		name           string
		id             string
		expectedResult bool
	}{
		{
			name:           "IsBanned",
			id:             usr.ID,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			isBanned, err := IsBanned(context.Background(), userService, tt.id)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedResult, isBanned)
		})
	}
}

func TestAddWarning(t *testing.T) {
	c, err := config.NewConfig()
	assert.NoError(t, err)

	client, err := InitClient(context.Background(), c.ConnURI)
	assert.NoError(t, err)

	defer func() {
		err = client.Disconnect(context.Background())
		assert.NoError(t, err)
	}()

	collection := mongodb.GetCollection(client, c.DB)
	userService := service.NewUserService(client, c.DB)

	usr, err := userService.SignUp(context.Background(), "testCli", "testcli@test.com", "testCli123")
	assert.NoError(t, err)

	defer func() {
		_, err := mongodb.DeleteAll(context.Background(), collection)
		assert.NoError(t, err)
	}()

	tests := []struct {
		name           string
		id             string
		expectedResult bool
		expectedError  error
	}{
		{
			name:           "AddWarning",
			id:             usr.ID,
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name:           "InvalidID",
			id:             usr.ID + "1",
			expectedResult: false,
			expectedError:  e.ErrDB,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			add, err := AddWarning(context.Background(), userService, tt.id)
			if err != nil {
				assert.Equal(t, tt.expectedError, err)
			}

			assert.Equal(t, tt.expectedResult, add)
		})
	}
}

func TestRemoveWarning(t *testing.T) {
	c, err := config.NewConfig()
	assert.NoError(t, err)

	client, err := InitClient(context.Background(), c.ConnURI)
	assert.NoError(t, err)

	defer func() {
		err = client.Disconnect(context.Background())
		assert.NoError(t, err)
	}()

	collection := mongodb.GetCollection(client, c.DB)
	userService := service.NewUserService(client, c.DB)

	usr, err := userService.SignUp(context.Background(), "testCli", "testcli@test.com", "testCli123")
	assert.NoError(t, err)

	defer func() {
		_, err := mongodb.DeleteAll(context.Background(), collection)
		assert.NoError(t, err)
	}()

	tests := []struct {
		name           string
		id             string
		expectedResult bool
	}{
		{
			name:           "AddWarning",
			id:             usr.ID,
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			_, err := AddWarning(context.Background(), userService, tt.id)
			assert.NoError(t, err)

			rem, err := RemoveWarning(context.Background(), userService, tt.id)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedResult, rem)
		})
	}
}

func TestRouter(t *testing.T) {
	c, err := config.NewConfig()
	assert.NoError(t, err)

	client, err := InitClient(context.Background(), c.ConnURI)
	assert.NoError(t, err)

	defer func() {
		err = client.Disconnect(context.Background())
		assert.NoError(t, err)
	}()

	collection := mongodb.GetCollection(client, c.DB)
	userService := service.NewUserService(client, c.DB)

	usr, err := userService.SignUp(context.Background(), "testCli", "testcli@test.com", "testCli123")
	assert.NoError(t, err)

	defer func() {
		_, err := mongodb.DeleteAll(context.Background(), collection)
		assert.NoError(t, err)
	}()

	tests := []struct {
		name           string
		id             string
		cmd            string
		expectedResult string
	}{
		{
			name:           "isBannedRouter",
			id:             usr.ID,
			cmd:            "isBan",
			expectedResult: "User with id: %v is not banned\n",
		},
		{
			name:           "AddWarningRouter",
			id:             usr.ID,
			cmd:            "add",
			expectedResult: "Add warning to user with id: %v\n",
		},
		{
			name:           "RemoveWarningRouter",
			id:             usr.ID,
			cmd:            "remove",
			expectedResult: "Remove warning for user with id: %v\n",
		},
		{
			name:           "CommandError",
			id:             usr.ID,
			cmd:            "",
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			str, err := Router(context.Background(), userService, tt.cmd, tt.id)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedResult, str)
		})
	}
}

func TestInitClient(t *testing.T) {
	c, err := config.NewConfig()
	assert.NoError(t, err)

	tests := []struct {
		name           string
		connURI        string
		expectedResult *mongo.Client
	}{
		{
			name:           "InitClient",
			connURI:        c.ConnURI,
			expectedResult: &mongo.Client{},
		},
		{
			name:           "InitClientError",
			connURI:        c.ConnURI + "1",
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			tClient, err := InitClient(context.Background(), c.ConnURI)
			assert.NoError(t, err)

			defer func() {
				err = tClient.Disconnect(context.Background())
				assert.NoError(t, err)
			}()

			assert.IsType(t, tClient, tt.expectedResult)
		})
	}
}
