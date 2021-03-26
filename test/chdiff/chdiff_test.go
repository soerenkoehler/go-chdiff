package chdiff_test

import (
	"path/filepath"
	"testing"

	"github.com/soerenkoehler/chdiff-go/main/chdiff"
	"github.com/soerenkoehler/chdiff-go/main/util"
	"github.com/soerenkoehler/go-testutils/mockutil"
)

type digestServiceMock struct {
	mockutil.Registry
}

func (mock digestServiceMock) Create(dataPath, digestPath, mode string) error {
	mockutil.Register(
		&mock.Registry,
		mockutil.Call{"create", dataPath, digestPath, mode})
	return nil
}

func (mock *digestServiceMock) Verify(dataPath, digestPath, mode string) error {
	mockutil.Register(
		&mock.Registry,
		mockutil.Call{"verify", dataPath, digestPath, mode})
	return nil
}

func expectDigestServiceCall(
	t *testing.T,
	args []string,
	call, dataPath, digestPath, mode string) {

	absDataPath, _ := filepath.Abs(dataPath)

	digestService := &digestServiceMock{}

	chdiff.DoMain(args, digestService, util.DefaultStdIOService{})

	mockutil.Verify(t,
		&digestService.Registry,
		mockutil.Call{call, absDataPath, digestPath, mode})
}

func TestCmdVerifyIsDefault(t *testing.T) {
	expectDigestServiceCall(t,
		[]string{""},
		"verify",
		".",
		"out.txt",
		"SHA256")
}

func TestCmdVerifyWithoutPath(t *testing.T) {
	expectDigestServiceCall(t,
		[]string{"", "v"},
		"verify",
		".",
		"out.txt",
		"SHA256")
}

func TestCmdVerifyWithPath(t *testing.T) {
	expectDigestServiceCall(t,
		[]string{"", "v", "x"},
		"verify",
		"x",
		"out.txt",
		"SHA256")
}
