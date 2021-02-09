package controller

import (
	"context"

	"github.com/fedo3nik/GamePortal_IdentityService/internal/application/service"
)

type CLIIsBannedHandler struct {
	userService service.User
}

type CLIAddWarningHandler struct {
	userService service.User
}

type CLIRemWarningHandler struct {
	userService service.User
}

func NewCLIIsBannedHandler(userService service.User) *CLIIsBannedHandler {
	return &CLIIsBannedHandler{userService: userService}
}

func (ch CLIIsBannedHandler) ServeCLI(ctx context.Context, id string) (bool, error) {
	_, isBanned, err := ch.userService.IsBanned(ctx, id)
	if err != nil {
		return false, err
	}

	return isBanned, err
}

func NewCLIAddWarningHandler(userService service.User) *CLIAddWarningHandler {
	return &CLIAddWarningHandler{userService: userService}
}

func (ch CLIAddWarningHandler) ServeCLI(ctx context.Context, id string) error {
	_, err := ch.userService.AddWarning(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func NewCLIRemWarningHandler(userService service.User) *CLIRemWarningHandler {
	return &CLIRemWarningHandler{userService: userService}
}

func (ch CLIRemWarningHandler) ServeCLI(ctx context.Context, id string) error {
	_, err := ch.userService.RemoveWarning(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
