package tokenize_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ikawaha/kagome/v2/cmd/tokenize"
)

func TestPrintScannedTokens_Default(t *testing.T) {
	userInput := "私"
	userArgs := []string{}

	// Mock STDIN
	if funcDefer, err := mockStdin(t, userInput); err != nil {
		t.Fatal(err)
	} else {
		defer funcDefer()
	}

	// Caputre output
	capturedSTDOUT := ""
	funcDefer := setCapturer(t, &capturedSTDOUT)

	defer funcDefer()

	// Run
	tokenize.Run(userArgs)

	// Assert
	expect := "私	名詞,代名詞,一般,*,*,*,私,ワタシ,ワタシ\nEOS\n"
	actual := capturedSTDOUT

	if expect != actual {
		t.Errorf("Expect: %v\nActual: %v", expect, actual)
	}
}

func TestPrintScannedTokens_JSON(t *testing.T) {
	userInput := "私"
	userArgs := []string{"-json"}

	if funcDefer, err := mockStdin(t, userInput); err != nil {
		t.Fatal(err)
	} else {
		defer funcDefer()
	}

	// Caputre output
	capturedSTDOUT := ""
	funcDefer := setCapturer(t, &capturedSTDOUT)

	defer funcDefer()

	// Run
	tokenize.Run(userArgs)

	// Assert
	expect := "[\n{\"id\":304999,\"start\":0,\"end\":1,\"surface\":\"私\"," +
		"\"class\":\"KNOWN\",\"pos\":[\"名詞\",\"代名詞\",\"一般\",\"*\"]," +
		"\"base_form\":\"私\",\"reading\":\"ワタシ\",\"pronunciation\":\"ワタシ\"," +
		"\"features\":[\"名詞\",\"代名詞\",\"一般\",\"*\",\"*\",\"*\",\"私\",\"ワタシ\"," +
		"\"ワタシ\"]}\n]\n"
	actual := capturedSTDOUT

	if expect != actual {
		t.Errorf("Expect: %v\nActual: %v", expect, actual)
	}
}

// Helper functions

// setCapturer is a helper function that captures the output of tokenize.FmtPrintF to capturedSTDOUT.
func setCapturer(t *testing.T, capturedSTDOUT *string) (funcDefer func()) {
	t.Helper()

	// Backup and set mock function
	oldFmtPrintF := tokenize.FmtPrintF
	tokenize.FmtPrintF = func(format string, a ...interface{}) (n int, err error) {
		*capturedSTDOUT += fmt.Sprintf(format, a...)

		return
	}

	// Return restore function
	return func() {
		tokenize.FmtPrintF = oldFmtPrintF
	}
}

// mockStdin is a helper function that lets the test pretend dummyInput as "os.Stdin" input.
// It will return a function for `defer` to clean up after the test.
func mockStdin(t *testing.T, dummyInput string) (funcDefer func(), err error) {
	t.Helper()

	oldOsStdin := os.Stdin
	tmpfile, err := ioutil.TempFile(t.TempDir(), t.Name())

	if err != nil {
		return nil, err
	}

	content := []byte(dummyInput)

	if _, err := tmpfile.Write(content); err != nil {
		return nil, err
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		return nil, err
	}

	// Set stdin to the temp file
	os.Stdin = tmpfile

	return func() {
		// clean up
		os.Stdin = oldOsStdin
		os.Remove(tmpfile.Name())
	}, nil
}
