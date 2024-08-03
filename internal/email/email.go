package email

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel"
)

const (
	// TODO(kg): Remove this.
	TestSubject = "You have requested to create an identity."
	TestBody    = "Do the following step to verify your identity."
	TestFrom    = "fugaziindustries@proton.me"
)

var tracer = otel.Tracer("email")

type Email struct {
	To      string
	Subject string
	Body    string
	From    string
}

type Sender interface {
	SendEmail(ctx context.Context, email Email) error
}

// NoOpSender is a no-op implementation of the Sender interface.
// It logs a fake email send to the console.
type NoOpSender struct{}

func NewNoOpSender() *NoOpSender {
	return &NoOpSender{}
}

// SendEmail logs the email to the console.
func (x *NoOpSender) SendEmail(ctx context.Context, email Email) error {
	ctx, span := tracer.Start(ctx, "SendEmail")
	defer span.End()

	slog.InfoContext(
		ctx,
		"no-op mailer: sending email",
		"to", email.To,
		"subject", email.Subject,
		"body", email.Body,
		"from", email.From,
	)
	return nil
}
