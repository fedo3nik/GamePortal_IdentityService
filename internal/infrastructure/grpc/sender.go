package grpc

import (
	"context"
	"io/ioutil"

	"github.com/fedo3nik/GamePortal_IdentityService/config"
)

type ServerGrpc struct {
	config *config.Config
}

func NewServerGrpc(c *config.Config) *ServerGrpc {
	return &ServerGrpc{config: c}
}

func (s *ServerGrpc) Send(context.Context, *Empty) (*SendResponse, error) {
	accessBytes, err := ioutil.ReadFile(s.config.AccessPublicKey)
	if err != nil {
		return nil, err
	}

	refreshBytes, err := ioutil.ReadFile(s.config.RefreshPublicKey)
	if err != nil {
		return nil, err
	}

	return &SendResponse{AccessPublicKey: string(accessBytes), RefreshPublicKey: string(refreshBytes)}, nil
}
