package generator

import (
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
)

func (g *Generator) generatePCSC() {
	var (
		constSource string
		constDest   string
		includePath string
	)

	switch g.os {
	case "linux":
		includePath = "/usr/include/PCSC/"
		constSource = filepath.Join(g.wd, "pcsc", "const_gen.go")
		constDest = filepath.Join(g.wd, "pcsc", "const.go")
	default:
		Fatal("pcsc unsupported os")
	}
	slog.Info("generating code for", slog.String("pkg", "github.com/happy-sdk/nfcsdk/pcsc"))

	cmd := exec.Command("go", "tool", "cgo", "-godefs", "--", "-I", includePath, constSource)
	cmd.Dir = g.temp
	constFile, err := os.Create(constDest)
	if err != nil {
		Fatalf("Failed to create output file: %w", err)
	}

	cmd.Stdout = constFile

	if err := cmd.Run(); err != nil {
		Fatalf("generatePCSC command failed: %w", err)
	}

	defer constFile.Close()
}
