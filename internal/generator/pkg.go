// Copyright 2023 The Happy Authors
// Licensed under the Apache License, Version 2.0.
// See the LICENSE file.

package generator

import "log/slog"

func (g *Generator) pkg() {
	slog.Info("generating code for", slog.String("pkg", "github.com/happy-sdk/nfcsdk"))
}
