package email

import (
	"bytes"
	"context"
	"encoding/gob"
	"log/slog"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("email")

const (
	IngestedEvent = "email_ingested"

	// TODO(kg): Remove this.
	TestSubject = "You have requested to create an identity."
	TestBody    = "Do the following step to verify your identity."
	TestFrom    = "fugaziindustries@proton.me"
)

type Email struct {
	To      string
	Subject string
	Body    string
	From    string
}

func (x Email) TraceAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("to", x.To),
		attribute.String("subject", x.Subject),
		attribute.String("Body", x.Body),
		attribute.String("From", x.From),
	}
}

func Ingest(ctx context.Context, natsConn *nats.Conn, email Email) error {
	ctx, span := tracer.Start(ctx, "Ingest", trace.WithAttributes(email.TraceAttributes()...))
	defer span.End()

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(email); err != nil {
		return err
	}
	if err := natsConn.Publish(IngestedEvent, buf.Bytes()); err != nil {
		return err
	}
	return nil
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
	ctx, span := tracer.Start(ctx, "SendEmail", trace.WithAttributes(email.TraceAttributes()...))
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
