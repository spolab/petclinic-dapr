package service

import (
	"context"
	"fmt"

	"github.com/dapr/go-sdk/client"
	"github.com/spolab/petclinic/owner/pkg/api/command"
)

func New(dapr client.Client, storename string) *OwnerService {
	if dapr == nil {
		panic("cannot use nil DAPR SDK")
	}
	return &OwnerService{dapr: dapr, storeName: storename}
}

type OwnerService struct {
	dapr      client.Client
	storeName string
}

func (s *OwnerService) Register(ctx context.Context, cmd *command.RegisterOwner) error {
	if cmd == nil {
		return fmt.Errorf("command cannot be nil")
	}
	// TODO Add validation
	s.dapr.GetState(ctx, s.storeName, cmd.GetOwner().GetId(), nil)
	err := s.dapr.SaveState(ctx, s.storeName, cmd.GetOwner().GetId(), []byte(cmd.String()), nil)
	return err
}
