package chdiff_test

import (
	"path/filepath"
	"testing"

	"github.com/soerenkoehler/chdiff-go/chdiff"
	"github.com/soerenkoehler/chdiff-go/util"
	"github.com/soerenkoehler/go-testutils/mockutil"
)

type digestServiceMock struct {
	mockutil.Registry
}

func (mock digestServiceMock) Create(dataPath, digestPath, algorithm string) {
	mockutil.Register(
		&mock.Registry,
		mockutil.Call{"create", dataPath, digestPath, algorithm})
}

func (mock *digestServiceMock) Verify(dataPath, digestPath, algorithm string) {
	mockutil.Register(
		&mock.Registry,
		mockutil.Call{"verify", dataPath, digestPath, algorithm})
}

func expectDigestServiceCall(
	t *testing.T,
	args []string,
	call, dataPath, digestPath, algorithm string) {

	absDataPath, _ := filepath.Abs(dataPath)
	absDigestPath, _ := filepath.Abs(digestPath)

	digestService := &digestServiceMock{}

	chdiff.DoMain("TEST", args, digestService, util.DefaultStdIOService{})

	mockutil.Verify(t,
		&digestService.Registry,
		mockutil.Call{call, absDataPath, absDigestPath, algorithm})
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
