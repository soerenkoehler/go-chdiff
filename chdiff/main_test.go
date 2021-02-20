package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/soerenkoehler/chdiff-go/util"
)

type digestServiceMock struct {
	calls []string
}

func (mock *digestServiceMock) Create(dataPath, digestPath, mode string) error {
	mock.calls = append(
		mock.calls,
		fmt.Sprintf("create %s %s %s", dataPath, digestPath, mode))
	return nil
}

func (mock *digestServiceMock) Verify(dataPath, digestPath, mode string) error {
	mock.calls = append(
		mock.calls,
		fmt.Sprintf("verify %s %s %s", dataPath, digestPath, mode))
	return nil
}

func expectDigestServiceCall(
	t *testing.T,
	args []string,
	call, dataPath, digestPath, mode string) {

	digestService := &digestServiceMock{}

	doMain(args, digestService, util.DefaultStdIOService{})

	absDataPath, _ := filepath.Abs(dataPath)
	expectedCalls := fmt.Sprintf(
		"%s %s %s %s",
		call,
		absDataPath,
		digestPath,
		mode)
	registeredCalls := strings.Join(digestService.calls, "\n")

	if expectedCalls != registeredCalls {
		t.Errorf(
			"\nexpected calls:\n%s\nregistered calls:\n%s",
			expectedCalls,
			registeredCalls)
	}
}

func TestCmdVerifyIsDefault(t *testing.T) {
	expectDigestServiceCall(
		t,
		[]string{""},
		"verify",
		".",
		"out.txt",
		"SHA256")
}

func TestCmdVerifyWithoutPath(t *testing.T) {
	expectDigestServiceCall(
		t,
		[]string{"", "v"},
		"verify",
		".",
		"out.txt",
		"SHA256")
}

func TestCmdVerifyWithPath(t *testing.T) {
	expectDigestServiceCall(
		t,
		[]string{"", "v", "x"},
		"verify",
		"x",
		"out.txt",
		"SHA256")
}
