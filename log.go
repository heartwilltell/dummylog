package dummylog

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/brianvoe/gofakeit"
)

const (
	// JSON represents text format for dummy logs.
	JSON Format = "json"
	// Text represents text format for dummy logs.
	Text Format = "text"

	// serverAddr holds default server address.
	serverAddr = ":8080"
	// serverShutdownTimeout holds default value for server shutdown timeout.
	serverShutdownTimeout = 3 * time.Second
)

// DummyLogger represents a process which will write
// silly messages to underlying writer.
// Implements io.Writer interface by writing silly
// message into given byte slice.
type DummyLogger struct {
	format  Format
	writer  io.Writer
	backoff time.Duration
	server  *http.Server
}

func New(options ...Option) *DummyLogger {
	handler := http.NewServeMux()

	l := DummyLogger{
		format:  Text,
		writer:  os.Stdout,
		backoff: time.Second,
		server: &http.Server{
			Addr:    serverAddr,
			Handler: handler,
		},
	}

	handler.HandleFunc("/say", l.say)

	for _, option := range options {
		option(&l)
	}

	return &l
}

func (l *DummyLogger) Write(p []byte) (int, error) {
	b := bytes.NewBuffer(p)

	if err := l.blabla(b); err != nil {
		return 0, err
	}

	return b.Len(), nil
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

func (l *DummyLogger) Serve(ctx context.Context) error {
	if l.server.Addr == "" {
		return fmt.Errorf("invalid server address")
	}

	// Handle shutdown signal in the background.
	go l.handleShutdown(ctx)

	if err := l.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTP server failed: %w", err)
	}

	return nil
}

func (l *DummyLogger) say(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	body, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(body); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// handleShutdown blocks until select statement receives a signal from
// ctx.Done, after that new context.WithTimeout will be created and passed to
// http.Server Shutdown method.
//
// If Shutdown method returns non nil error, program will call os.Exit immediately.
func (l *DummyLogger) handleShutdown(ctx context.Context) {
	<-ctx.Done()

	killctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()

	l.server.SetKeepAlivesEnabled(false)

	if err := l.server.Shutdown(killctx); err != nil {
		panic(fmt.Errorf("HTTP server force exit! Failed to shutdown the HTTP server gracefully: %w", err))
	}
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
