package notification

import "context"

type PushNotificationHandler struct {
	firebaseKey string
}

func (ph *PushNotificationHandler) Handle(ctx context.Context, msg []byte) error {
	return nil
}
