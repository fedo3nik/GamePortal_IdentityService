package main

import (
	"context"
	"log"
	"os"
	"time"

	cliController "github.com/fedo3nik/GamePortal_IdentityService/internal/interface/controller/cli"

	"github.com/fedo3nik/GamePortal_IdentityService/config"
	"github.com/fedo3nik/GamePortal_IdentityService/internal/application/service"
	"github.com/mkideal/cli"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const contextTimeout = 10

type argT struct {
	ID      string `cli:"id" usage:"enter the user's id who should be ban'"`
	Command string `cli:"cmd" usage:"enter command what you should to use"`
}

func InitClient(ctx context.Context, connURI string) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connURI))
	if err != nil {
		log.Printf("InitClient error: %v", err)
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Printf("Connect error: %v", err)
		return nil, err
	}

	return client, nil
}

func IsBanned(ctx context.Context, userService service.User, id string) (bool, error) {
	isBannedHandler := cliController.NewCLIIsBannedHandler(userService)

	isBanned, err := isBannedHandler.ServeCLI(ctx, id)
	if err != nil {
		return false, err
	}

	if isBanned {
		return true, err
	}

	return false, nil
}

func AddWarning(ctx context.Context, userService service.User, id string) (bool, error) {
	addWarningHandler := cliController.NewCLIAddWarningHandler(userService)

	err := addWarningHandler.ServeCLI(ctx, id)
	if err != nil {
		log.Print(err)
		return false, err
	}

	return true, nil
}

func RemoveWarning(ctx context.Context, userService service.User, id string) (bool, error) {
	remWarningHandler := cliController.NewCLIRemWarningHandler(userService)

	err := remWarningHandler.ServeCLI(ctx, id)
	if err != nil {
		log.Print(err)
		return false, err
	}

	return true, nil
}

func Router(ctx context.Context, userService service.User, cmd, id string) (string, error) {
	switch cmd {
	case "add":
		{
			_, err := AddWarning(ctx, userService, id)
			if err != nil {
				return "", err
			}

			return "Add warning to user with id: %v\n", nil
		}
	case "remove":
		{
			_, err := RemoveWarning(ctx, userService, id)
			if err != nil {
				return "", err
			}

			return "Remove warning for user with id: %v\n", nil
		}
	case "isBan":
		{
			isUserBanned, err := IsBanned(ctx, userService, id)
			if err != nil {
				return "", err
			}

			if isUserBanned {
				return "User with id: %v is banned\n", nil
			}

			return "User with id: %v is not banned\n", nil
		}
	}

	return "", nil
}

func main() {
	os.Exit(cli.Run(new(argT), func(cliCtx *cli.Context) error {
		c, err := config.NewConfig()
		if err != nil {
			log.Panic(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), contextTimeout*time.Second)

		client, err := InitClient(ctx, c.ConnURI)
		if err != nil {
			log.Panic(err)
		}

		defer func() {
			cancel()

			err = client.Disconnect(ctx)
			if err != nil {
				log.Panicf("Disconnecr error: %v", err)
			}
		}()

		userService := service.NewUserService(client, c.DB)

		argv := cliCtx.Argv().(*argT)

		routerMessage, err := Router(ctx, userService, argv.Command, argv.ID)
		if err != nil {
			log.Print(err)
		}

		cliCtx.String(routerMessage, argv.ID)

		return nil
	}))
}
