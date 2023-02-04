package handlers_test

import (
	"os"
	"receiver/logger"
	"testing"
)

//!+test
//go test -v

func TestMain(m *testing.M) {
	log := logger.NewEmptyLogger()
	logger.InitializeLogger(&log)
	code := m.Run()
	os.Exit(code)
}

//!-tests
