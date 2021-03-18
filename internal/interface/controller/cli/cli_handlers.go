package cli

import (
	"context"

	"github.com/fedo3nik/GamePortal_IdentityService/internal/application/service"
)

type IsBannedHandler struct {
	userService service.User
}

type AddWarningHandler struct {
	userService service.User
}

type RemWarningHandler struct {
	userService service.User
}

func NewCLIIsBannedHandler(userService service.User) *IsBannedHandler {
	return &IsBannedHandler{userService: userService}
}

func (ch IsBannedHandler) ServeCLI(ctx context.Context, id string) (bool, error) {
	_, isBanned, err := ch.userService.IsBanned(ctx, id)
	if err != nil {
		return false, err
	}

	return isBanned, err
}

func NewCLIAddWarningHandler(userService service.User) *AddWarningHandler {
	return &AddWarningHandler{userService: userService}
}

func (ch AddWarningHandler) ServeCLI(ctx context.Context, id string) error {
	_, err := ch.userService.AddWarning(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func NewCLIRemWarningHandler(userService service.User) *RemWarningHandler {
	return &RemWarningHandler{userService: userService}
}

func (ch RemWarningHandler) ServeCLI(ctx context.Context, id string) error {
	_, err := ch.userService.RemoveWarning(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
