package generator

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Fatal logs the provided message as an error and then exits the program with exit code 1.
func Fatal(msg string) {
	slog.Error(msg)
	os.Exit(1)
}

func Fatalf(format string, a ...any) {
	Fatal(fmt.Sprintf(format, a...))
}

// Fatalf formats the provided message according to a format specifier and then logs it as an error.
// It then exits the program with exit code 1. It is analogous to fmt.Printf, but for error logging.
func FileExists(filename string) error {
	info, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", filename)
		}
		return fmt.Errorf("error checking file: %s, %v", filename, err)
	}
	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a file: %s", filename)
	}
	return nil
}

type Generator struct {
	valid bool
	wd    string
	temp  string
	os    string
}

func NewGenerator() (*Generator, error) {
	slog.Info("preparing code generator")
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	temp, err := os.MkdirTemp("", "pcnfcsdk-")
	if err != nil {
		Fatal(fmt.Errorf("failed to create temporary directory: %w", err).Error())
	}

	g := &Generator{
		wd:   wd,
		temp: temp,
		os:   runtime.GOOS,
	}

	switch runtime.GOOS {
	case "linux":
	default:
		Fatalf("unsupprted os %s", runtime.GOOS)
	}

	if err := FileExists(filepath.Join(wd, "nfcsdk.go")); err != nil {
		return nil, err
	}

	slog.Info("generator ready", g.info())
	g.valid = true
	return g, nil
}

func (g *Generator) Finalize() {
	slog.Info("finalizing generator")
	// ALWAYS

	// remove tmp dir
	if len(g.temp) > 0 {
		os.RemoveAll(g.temp)
		slog.Info("removed", slog.String("temp", g.temp))
	}

	// ONLY ON VALID STATE
	if g.valid {
		fmtCmd := exec.Command("go", "fmt", "./...")
		fmtCmd.Dir = g.wd
		if err := fmtCmd.Run(); err != nil {
			Fatal(fmt.Sprintf("failed to format code: %w", err))
		}
		slog.Info("go fmt completed")
	}
}

func (g *Generator) Generate() {
	g.generatePKG()
	g.generateNDEF()
	g.generatePCSC()
}

func (g *Generator) info() slog.Attr {
	return slog.Group(
		"generator",
		slog.String("wd", g.wd),
	)
}
