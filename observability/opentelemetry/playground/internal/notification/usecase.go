package notification

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/tuanta7/monitor/pkg/monitor"
)

type UseCase struct {
	tracer *monitor.Tracer
}

func NewUseCase(tracer *monitor.Tracer) *UseCase {
	return &UseCase{
		tracer: tracer,
	}
}

func (u *UseCase) SendPushNotification(ctx context.Context) error {
	ctx, span := u.tracer.StartNewSpan(ctx, "notification.push")
	defer span.End()

	r := rand.Intn(10)
	if r > 9 {
		time.Sleep(100 * time.Millisecond)
		err := errors.New("failed to push notification")
		span.RecordError(err)
		return err
	}

	time.Sleep(30 * time.Millisecond)
	return nil
}
