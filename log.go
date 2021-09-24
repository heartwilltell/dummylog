package dummylog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/brianvoe/gofakeit"
)

const (
	// JSON represents text format for dummy logs.
	JSON Format = "json"
	// Text represents text format for dummy logs.
	Text Format = "text"
)

// DummyLogger represents a process which will write
// silly messages to underlying writer.
// Implements io.Writer interface by writing silly
// message into given byte slice.
type DummyLogger struct {
	format  Format
	writer  io.Writer
	backoff time.Duration
}

func New(options ...Option) *DummyLogger {
	l := DummyLogger{
		format:  Text,
		writer:  os.Stdout,
		backoff: time.Second,
	}

	for _, option := range options {
		option(&l)
	}

	return &l
}

func (l *DummyLogger) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(l.backoff):
			if err := l.blabla(l.writer); err != nil {
				return err
			}

			continue
		}
	}
}

func (l *DummyLogger) Write(p []byte) (int, error) {
	b := bytes.NewBuffer(p)

	if err := l.blabla(b); err != nil {
		return 0, err
	}

	return b.Len(), nil
}

func (l *DummyLogger) blabla(w io.Writer) error {
	switch l.format {
	case JSON:
		message := fmt.Sprintf(`{"message": "%s"}`, gofakeit.HipsterSentence(5)) + "\n"
		if _, err := fmt.Fprintln(w, message); err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}

	case Text:
		if _, err := fmt.Fprintf(w, "%s\n", gofakeit.HipsterSentence(5)); err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}
	}

	return nil
}

// Option represents functional options pattern for DummyLogger type,
// a function which receive a pointer to DummyLogger struct.
// Option functions can only be passed to DummyLogger constructor function New,
// and can change the default value of DummyLogger struct.
type Option func(l *DummyLogger)

// WithWriter sets DummyLogger writer.
func WithWriter(writer io.Writer) Option { return func(l *DummyLogger) { l.writer = writer } }

// WithFormat sets DummyLogger log format.
func WithFormat(format Format) Option { return func(l *DummyLogger) { l.format = format } }

type Format string

func (f Format) String() string { return string(f) }
