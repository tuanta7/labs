package main

import (
	"context"
)

type SubscribeHandler struct{}

func (h *SubscribeHandler) Handle(ctx context.Context) error {
	return nil
}
