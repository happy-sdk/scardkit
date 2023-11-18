package generator

import "log/slog"

func (g *Generator) generateNDEF() {
	slog.Info("generating code for", slog.String("pkg", "github.com/happy-sdk/nfcsdk/ndef"))
}
