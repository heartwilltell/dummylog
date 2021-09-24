package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/heartwilltell/dummylog"
)

const (
	runCmd = "run"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			os.Exit(1)
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if len(os.Args) <= 1 {
		panic(fmt.Errorf("dummylog: expected one of the following commands: '%s'", runCmd))
	}

	switch os.Args[1] {
	case runCmd:
		if err := runCommand(ctx, os.Args[1:]); err != nil && !errors.Is(err, context.Canceled) {
			panic(fmt.Errorf("dummylog: %s command failed: %w", runCmd, err))
		}
	default:
		panic(fmt.Errorf("dummylog: unexpected command: '%s', use '%s' to run log writer", os.Args[1], runCmd))
	}
}

func runCommand(ctx context.Context, args []string) error {
	var (
		format   string
		filePath string
	)

	cmd := flag.NewFlagSet(runCmd, flag.ExitOnError)
	cmd.StringVar(&format, "format", "text", "sets log format: 'text' or 'json'.")
	cmd.StringVar(&filePath, "file", "", "sets path to file where logs will be written.")

	if err := cmd.Parse(args[1:]); err != nil {
		return fmt.Errorf("failed to parse command arguments: %w", err)
	}

	options := make([]dummylog.Option, 0, 2)

	if filePath != "" {
		file, createErr := os.Create(filePath)
		if createErr != nil {
			return createErr
		}

		options = append(options, dummylog.WithWriter(file))
	}

	switch format {
	case "json":
		options = append(options, dummylog.WithFormat(dummylog.JSON))
	case "text":
		options = append(options, dummylog.WithFormat(dummylog.Text))
	default:
		return fmt.Errorf("unknown format: %s", format)
	}

	dummy := dummylog.New(options...)
	if err := dummy.Start(ctx); err != nil {
		return err
	}

	return nil
}
