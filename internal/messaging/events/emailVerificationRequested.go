package events

type EmailVerificationRequested struct {
	Event
	Email string `json:"email"`
}
