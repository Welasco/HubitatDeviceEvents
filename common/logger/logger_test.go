package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitWithoutFileLogsToStdout(t *testing.T) {
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	original := os.Stdout
	os.Stdout = writer
	defer func() {
		os.Stdout = original
		l.SetOutput(original)
		writer.Close()
		reader.Close()
	}()
	if err := Init("", "Inf"); err != nil {
		t.Fatal(err)
	}
	Info("[test][stdout] shutdown message")
	writer.Close()
	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(output), "[test][stdout] shutdown message") {
		t.Fatalf("stdout is missing log message: %s", output)
	}
}

func TestInitReportsFileOpenFailure(t *testing.T) {
	err := Init(filepath.Join(t.TempDir(), "missing-directory", "events.log"), "Inf")
	if err == nil {
		t.Fatal("expected explicit error opening optional log file")
	}
}
