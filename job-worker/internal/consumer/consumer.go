package consumer

import (
	"context"
	"fmt"
)

var ErrConsumerClosed = fmt.Errorf("consumer closed")

type Consumer interface {
	Consume() error
}

type consumeCallback[T any] func(context.Context, T, error) error
