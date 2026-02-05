package notification

type UseCase struct {
}

func NewUseCase() *UseCase {
	return &UseCase{}
}

func (uc *UseCase) SendEmail(message string) error {
	return nil
}

func (uc *UseCase) SendPushNotification(message string) error {
	return nil
}
