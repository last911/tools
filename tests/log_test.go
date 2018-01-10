package tests

import (
	"github.com/last911/tools/log"
	"testing"
)

func TestLog(t *testing.T) {
	log.InitLogger("")

	log.Info("Test log format")
}
