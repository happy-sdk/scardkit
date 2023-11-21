// Copyright 2023 The Happy Authors
// Licensed under the Apache License, Version 2.0.
// See the LICENSE file.

// Package scardkit offers a robust suite for smart card interactions,
// integrating various protocols, PC/SC and NFC technologies.
// It's designed to abstract the complexities involved in smart card
// communications, providing developers with a user-friendly and consistent API.
// This toolkit caters to a wide range of smart card applications, from security
// authentication to data management, and is crafted with the idiomatic Go practices
// and simplicity in mind.
package scardkit

// New initializes a new instance of the smart card SDK.
func New() *SDK { return nil }

// Run executes a provided Command and returns a Response.
func (sdk *SDK) Run(cmd Command) (Response, error) { return nil, nil }

// SDK represents the smart card toolkit with common functionalities.
type SDK struct {
	// Fields for SDK configuration and state
}

// Command represents a generic command interface that can be implemented by different card protocols.
type Command interface {
	Execute() (Response, error)
}

// Response is a generic interface for responses from smart card operations.
type Response interface {
	// Methods to process and retrieve response data
}

// Feature represents the state of a feature in the scardkit subpackages.
// It is used to define the availability and status of features based on the platform.
type Feature uint8

const (
	FeatureNotImplemented Feature = iota // Indicates the feature is not implemented in the package.
	FeatureImplemented                   // Indicates the feature is implemented but its status (enabled/disabled) is not specified.
	FeatureEnabled                       // Indicates the feature is implemented and currently enabled.
	FeatureDisabled                      // Indicates the feature is implemented but currently disabled.
)
