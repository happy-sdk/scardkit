package generator

import "log/slog"

func (g *Generator) generatePKG() {
	slog.Info("generating code for", slog.String("pkg", "github.com/happy-sdk/nfcsdk"))
}
