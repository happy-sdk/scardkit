// Copyright 2023 The Happy Authors
// Licensed under the Apache License, Version 2.0.
// See the LICENSE file.

// Package nfcsdk provides a comprehensive software development kit for
// integrating and manipulating NTAG213, NTAG215, and NTAG216 NFC tags.
// It includes robust tools for reading and writing operations, enabling
// efficient handling of NFC tag data in various applications. This SDK
// is optimized for developers seeking a reliable and streamlined solution
// for NFC tag interactions.
package nfcsdk

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"
)

//go:generate go run internal/generator/main.go

var (
	Error             = errors.New("error")
	ErrInvalidContext = fmt.Errorf("%w: a nil or cancellable context is required, context.TODO or context.Background are not valid", Error)
)

func New(ctx context.Context, logger *slog.Logger) (*SDK, error) {
	sdk := &SDK{}

	// set log/slog if present
	if logger != nil {
		sdk.logger = logger.WithGroup("nfc")
	}

	// use provided context
	if ctx != nil {
		if c, ok := ctx.(fmt.Stringer); ok {
			switch c.String() {
			case "context.TODO", "context.Background":
				return nil, fmt.Errorf("%w: ", ErrInvalidContext)
			}
		}
		sdk.ctx, sdk.stop = context.WithCancelCause(ctx)

	} else {
		var stop context.CancelFunc
		sdk.ctx, stop = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		sdk.stop = func(cause error) {
			if cause != nil {
				sdk.error(cause)
			}
			stop()
		}
	}

	if err := sdk.init(); err != nil {
		sdk.error(err)
		sdk.dispose()
		return nil, err
	}
	sdk.wg.Add(1)
	go func() {
		defer sdk.stop(nil)
		<-sdk.ctx.Done()
		if err := sdk.ctx.Err(); err != nil {
			sdk.warn(err.Error())
		}
		sdk.dispose()
		sdk.wg.Done()
	}()
	return sdk, nil
}

type Reader struct {
	id   int
	name string
	Use  bool // When set truie then this reader listens for cards
}

func (r Reader) Name() string {
	return r.name
}

func (r Reader) ID() int {
	return r.id
}

type ReaderSelectFunc func(rs []Reader) (use []Reader, err error)
