// Copyright 2023 The Happy Authors
// Licensed under the Apache License, Version 2.0.
// See the LICENSE file.

package nfcsdk

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"sync"
	"time"

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
	cardHandler  CardHandler

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
		return
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
		return
	}

	s.mu.RLock()
	readers := s.readers
	s.mu.RUnlock()
	var states []pcsc.ReaderState
	for _, reader := range readers {
		if !reader.Use {
			continue
		}
		states = append(states, pcsc.ReaderState{
			Reader:       reader.name, // Replace with the actual reader name
			CurrentState: pcsc.ScardStateUnaware,
		})
	}

	if len(states) == 0 {
		err = fmt.Errorf("%w: no readers enabled", Error)
		s.stop(err)
		s.wg.Wait()
		return
	}
	s.info("started", slog.String("time", time.Now().String()))
runner:
	for {
		select {
		case <-s.ctx.Done():
			break runner
		default:
			// check is context valid
			if err = s.hctx.IsValid(); err != nil {
				s.error(err)
				break runner
			}
			err = s.hctx.GetStatusChange(states, -1)
			if err != nil {
				s.error(err)
				break runner
			}

			for i := range states {
				states[i].CurrentState = states[i].EventState
				if states[i].EventState&pcsc.ScardStatePresent != 0 {
					s.debug("card is present")
					// check again context mat get invalid
					if err = s.hctx.IsValid(); err != nil {
						s.error(err)
						break runner
					}

					s.handleCard(states[i].Reader)

				} else {
					s.debug("no card present, waiting...")
				}
			}

		}

	}

	s.wg.Wait() // Wait for shutdown and cleanup

	s.debug("exiting")
	return
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

func (s *SDK) OnCardPresent(fn CardHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cardHandler != nil {
		s.warn("reader select callback can only be attached once")
		return
	}
	s.cardHandler = fn
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
		s.debug("found reader", slog.Group("reader", slog.Int("id", reader.id), slog.String("name", readerName)))
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
	s.debug("disposing...")
	if s.hctx != nil {
		s.debug("cancel pending actions")
		if err := s.hctx.Cancel(); err != nil {
			s.error(err)
		}
		s.debug("released scard context")
		if err := s.hctx.Release(); err != nil {
			s.error(err)
		}
	}

	s.debug("sdk disposed")
}

func (s *SDK) handleCard(readerName string) error {
	card, err := s.hctx.Connect(readerName, pcsc.ScardShareExclusive, pcsc.ScardProtocolAny)
	if err != nil {
		s.error(err)
		return nil
	}
	s.debug("card connected", slog.String("protocols", card.Protocol().String()))

	if err := card.RefreshStatus(); err != nil {
		s.error(err)
		return nil
	}
	status := card.CurrentStatus()
	s.info("connected card status",
		slog.String("state", status.State.String()),
		slog.String("protocol", status.Protocol.String()),
		slog.String("reader", status.Reader),
		slog.String("atr", FormatByteSlice(status.Atr)),
	)

	// Call user callback
	s.mu.RLock()
	handler := s.cardHandler
	s.mu.RUnlock()
	if handler != nil {
		if err := handler(card); err != nil {
			s.error(err)
			// Disconnect from card
			if err2 := card.Disconnect(pcsc.ScardResetCard); err2 != nil {
				s.error(err2)
				return errors.Join(err, err2)
			}
			s.info("card disconnected")
			return err
		}
	} else {
		s.info("no card handler defined")
	}
	// Disconnect from card
	if err := card.Disconnect(pcsc.ScardResetCard); err != nil {
		s.error(err)
		return nil
	}

	s.info("card disconnected")
	return nil
}

const logPrefix = "nfc: "

// LogAttrs is a more efficient version of [Logger.Log] that accepts only Attrs.
func (s *SDK) log(level slog.Level, msg string, args ...any) {
	if s.logger == nil {
		return
	}
	if !s.logger.Enabled(context.Background(), level) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // skip [Callers, Infof]
	msg = logPrefix + msg
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(args...)
	_ = s.logger.Handler().Handle(context.Background(), r)
}

func (s *SDK) debug(msg string, args ...any) {
	s.log(slog.LevelDebug, msg, args...)
}

func (s *SDK) info(msg string, args ...any) {
	s.log(slog.LevelInfo, msg, args...)
}

func (s *SDK) warn(msg string, args ...any) {
	s.log(slog.LevelWarn, msg, args...)
}

func (s *SDK) error(err error) {
	if err == nil {
		return
	}
	s.log(slog.LevelError, err.Error())
}
