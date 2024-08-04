package event

import (
	"bytes"
	"context"
	"encoding/gob"
	"log/slog"

	"github.com/Salam4nder/identity/internal/email"
	"github.com/nats-io/nats.go"
)

type Worker struct {
	mailSender email.Sender
}

func NewWorker(sender email.Sender) *Worker {
	return &Worker{
		mailSender: sender,
	}
}

func (x *Worker) Work(ctx context.Context, natsCh chan *nats.Msg) {
	for {
		select {
		case <-ctx.Done():
			slog.InfoContext(ctx, "event: context done, shutting down worker...")
			return
		case msg := <-natsCh:
			var m email.Email
			if err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(&m); err != nil {
				slog.WarnContext(ctx, "event: decoding message", "err", err)
				continue
			}
			if err := x.mailSender.SendEmail(ctx, email.Email{
				To:      m.To,
				Subject: m.Subject,
				Body:    m.Body,
				From:    m.From,
			}); err != nil {
				slog.WarnContext(ctx, "event: sending email", "err", err)
			}
		}
	}
}
