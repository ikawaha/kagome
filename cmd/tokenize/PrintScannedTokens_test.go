package tokenize_test

import (
	"errors"
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
	err := tokenize.Run(userArgs)
	if err != nil {
		t.Fatalf("Failed to execute tokenize.Run.\n%v", err)
	}

	// Assert
	actual := capturedSTDOUT
	expect := "私	名詞,代名詞,一般,*,*,*,私,ワタシ,ワタシ\nEOS\n"
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
	err := tokenize.Run(userArgs)
	if err != nil {
		t.Fatalf("Failed to execute tokenize.Run.\n%v", err)
	}

	// Assert
	actual := capturedSTDOUT
	expect := "[\n{\"id\":304999,\"start\":0,\"end\":1,\"surface\":\"私\"," +
		"\"class\":\"KNOWN\",\"pos\":[\"名詞\",\"代名詞\",\"一般\",\"*\"]," +
		"\"base_form\":\"私\",\"reading\":\"ワタシ\",\"pronunciation\":\"ワタシ\"," +
		"\"features\":[\"名詞\",\"代名詞\",\"一般\",\"*\",\"*\",\"*\",\"私\",\"ワタシ\"," +
		"\"ワタシ\"]}\n]\n"

	if expect != actual {
		t.Errorf("Expect: %v\nActual: %v", expect, actual)
	}
}

// TestPrintScannedTokens_parse_fail covers the json.Marshal failure.
func TestPrintScannedTokens_parse_fail(t *testing.T) {
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

	// Backup JSONMarshal and restore
	oldJSONMarshal := tokenize.JSONMarshal
	defer func() {
		tokenize.JSONMarshal = oldJSONMarshal
	}()

	// Mock JSONMarshal
	msgError := "forced fail"
	tokenize.JSONMarshal = func(v interface{}) ([]byte, error) {
		return nil, errors.New(msgError)
	}

	// Run
	err := tokenize.Run(userArgs)
	if err == nil {
		t.Fatalf("failure test failed. The tokenize.Run should return an error")
	}

	// Assert
	expect := msgError
	actual := err.Error()

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
		err := os.Remove(tmpfile.Name())
		if err != nil {
			t.Fatalf("failed to remove temp file during test.\n%v", err)
		}
	}, nil
}
