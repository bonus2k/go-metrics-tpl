package osexit_test

import (
	"github.com/bonus2k/go-metrics-tpl/internal/analyzer/osexit"
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

func TestOsExit(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), osexit.Analyzer, "./...")
}
