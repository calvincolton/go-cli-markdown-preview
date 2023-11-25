package main

import (
	"bytes"
	"os"
	"regexp"
	"strings"
	"testing"
)

const (
	inputFile  = "./testdata/test1.md"
	goldenFile = "./testdata/test1.md.html"
)

func normalizeContent(content []byte) string {
	re := regexp.MustCompile(`[\s\t\n]+`)
	return string(re.ReplaceAll(content, []byte(" ")))
}

func TestParseContent(t *testing.T) {
	input, err := os.ReadFile(inputFile)
	if err != nil {
		t.Fatal(err)
	}

	result := parseContent(input)

	expected, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}

	expectedNormalized := normalizeContent(expected)
	resultNormalized := normalizeContent(result)

	if expectedNormalized != resultNormalized {
		t.Logf("golden:\n%s\n", expectedNormalized)
		t.Logf("result:\n%s\n", resultNormalized)
		t.Error("Result content does not match golden file")
	}
}

func TestRun(t *testing.T) {
	var mockStdOut bytes.Buffer

	if err := run(inputFile, &mockStdOut); err != nil {
		t.Fatal(err)
	}
	resultFile := strings.TrimSpace(mockStdOut.String())

	result, err := os.ReadFile(resultFile)
	if err != nil {
		t.Fatal(err)
	}

	expected, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}

	expectedNormalized := normalizeContent(expected)
	resultNormalized := normalizeContent(result)

	if expectedNormalized != resultNormalized {
		t.Logf("golden:\n%s\n", expectedNormalized)
		t.Logf("result:\n%s\n", resultNormalized)
		t.Error("Result content does not match golden file")
	}

	os.Remove(resultFile)
}
