package main

import (
	"github.com/Khasmag06/go-url-shortener/cmd/staticlint/analyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

func main() {
	// определяем map подключаемых правил
	checks := map[string]bool{
		"ST1000": true,
		"ST1005": true,
		"ST1023": true,
	}

	myChecks := []*analysis.Analyzer{
		analyzer.OsExitAnalyzer,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
	}
	for _, v := range staticcheck.Analyzers {
		myChecks = append(myChecks, v.Analyzer)
	}

	for _, v := range stylecheck.Analyzers {
		if checks[v.Analyzer.Name] {
			myChecks = append(myChecks, v.Analyzer)
		}
	}

	multichecker.Main(
		myChecks...,
	)
}
