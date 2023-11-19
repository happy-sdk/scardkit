// Copyright 2023 The Happy Authors
// Licensed under the Apache License, Version 2.0.
// See the LICENSE file.

package nfcsdk

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/happy-sdk/nfcsdk/pcsc"
)

type SDK struct {
	mu           sync.RWMutex
	ctx          context.Context
	stop         context.CancelCauseFunc
	logger       *slog.Logger
	disposed     bool
	wg           sync.WaitGroup
	readerSelect ReaderSelectFunc

	hctx    *pcsc.HContext
	readers []Reader
}

func (s *SDK) Disposed() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	disposed := s.disposed
	return disposed
}

// Run initiates the SDK's main operations and then blocks until the SDK's context has been cancelled.
// After cancellation, it ensures that all cleanup operations within the SDK are fully completed before returning.
// This function is used to both start the SDK and to wait for its graceful shutdown,
// making sure that all resources are properly released and all operations are concluded.
// If an error occurs during execution, the error is returned after all resources are freed and cleanup operations are completed,
// ensuring that the SDK is always left in a clean state regardless of whether it terminated successfully or due to an error.
func (s *SDK) Run() (err error) {
	s.mu.RLock()
	// Context can get invalidated.
	if err = s.hctx.IsValid(); err != nil {
		s.error(err)
	}
	if len(s.readers) == 0 {
		err = errors.Join(err, fmt.Errorf("%w: no readers present", Error))
	}
	s.mu.RUnlock()
	if err != nil {
		s.stop(err)
		s.wg.Wait()
		return err
	}

	// Select readers
	s.mu.Lock()
	if s.readerSelect != nil {
		readers := s.readers
		s.readers, err = s.readerSelect(readers)
	} else {
		s.readers[0].Use = true // use by default only first reader
	}
	s.mu.Unlock()
	if err != nil {
		s.stop(err)
		s.wg.Wait()
		return err
	}

	s.wg.Wait() // Wait for shutdown and cleanup
	return nil
}

// SelectReader allows for specifying a callback function (fn) that determines the selection
// criteria for the reader. This callback is used to choose which reader to use based on
// custom logic provided in fn. If no callback is set (fn is nil), the SDK defaults to
// selecting the first available reader. Note that the callback can only be set once;
// attempting to set it again will result in a warning and the subsequent call will be ignored.
func (s *SDK) SelectReader(fn ReaderSelectFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.readerSelect != nil {
		s.warn("reader select callback can only be attached once")
		return
	}

	s.readerSelect = fn
}

func (s *SDK) init() (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// call SCardEstablishContext
	s.hctx, err = pcsc.SCardEstablishContext(pcsc.ScardScopeSystem)
	if err != nil {
		return err
	}
	s.debug("scard context established")

	// Get available readers
	readerNames, err := s.hctx.ListReaders()
	if err != nil {
		return err
	}
	if len(readerNames) == 0 {
		return fmt.Errorf("%w: no readers returned by ListReaders", Error)
	}
	for i, readerName := range readerNames {
		reader := Reader{
			id:   i + 1,
			name: readerName,
		}
		s.debug("found", slog.Group("reader", slog.Int("id", reader.id), slog.String("name", readerName)))
		s.readers = append(s.readers, reader)
	}

	return nil
}

func (s *SDK) dispose() {
	if s.Disposed() {
		s.warn("sdk already disposed")
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.disposed = true
	s.debug("disposing")
	if s.hctx != nil {
		if err := s.hctx.Release(); err != nil {
			s.error(err)
		} else {
			s.debug("scard context released")
		}
	}

	s.debug("sdk disposed")
}

const logPrefix = "nfc: "

// LogAttrs is a more efficient version of [Logger.Log] that accepts only Attrs.
func (s *SDK) Log(level slog.Level, msg string, args ...any) {
	if s.logger == nil {
		return
	}
	msg = logPrefix + msg
	s.logger.Log(s.ctx, level, msg, args...)
}

func (s *SDK) debug(msg string, args ...any) {
	s.Log(slog.LevelDebug, msg, args...)
}

func (s *SDK) info(msg string, args ...any) {
	s.Log(slog.LevelInfo, msg, args...)
}

func (s *SDK) warn(msg string, args ...any) {
	s.Log(slog.LevelWarn, msg, args...)
}

func (s *SDK) error(err error) {
	if err == nil {
		return
	}
	s.Log(slog.LevelError, err.Error())
}
