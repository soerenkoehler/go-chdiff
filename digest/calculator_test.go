package digest_test

import (
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/soerenkoehler/go-chdiff/common"
	"github.com/soerenkoehler/go-chdiff/digest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type testCase struct {
	path     string
	size     int64
	seed     int64
	hash     string
	inDigest bool
}

type TestSuiteCalculator struct {
	suite.Suite
	root string
}

func TestSuiteRunner(t *testing.T) {
	suite.Run(t, &TestSuiteCalculator{})
}

func (s *TestSuiteCalculator) SetupTest() {
	s.root = s.T().TempDir()
}

func (s *TestSuiteCalculator) TestDigest256() {
	s.verifyDigest([]testCase{{
		path:     "zero",
		size:     0,
		seed:     1,
		hash:     "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		inDigest: true,
	}, {
		path:     "data_1",
		size:     256,
		seed:     1,
		hash:     "352adfeb0dc6e28699635c5911cf33e2e0a86aedf85a5a99bba97749000ae1c7",
		inDigest: true,
	}, {
		path:     "sub/data_1",
		size:     256,
		seed:     2,
		hash:     "1629705c76a590f2e16b8c42fa0aca9c405401fcfc794399e71f0954f1e0d50e",
		inDigest: true,
	}}, digest.SHA256)
}

func (s *TestSuiteCalculator) TestDigest512() {
	s.verifyDigest([]testCase{{
		path:     "zero",
		size:     0,
		seed:     1,
		hash:     "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e",
		inDigest: true,
	}, {
		path:     "data_1",
		size:     256,
		seed:     1,
		hash:     "aa43d14bc209ae859af792d9d0ba6ab27ab7d3802281c6a528485d44ac18f1c5019287e93ec1d3f15e843df0f05278b06471e61597b05cee6d3a347434729b88",
		inDigest: true,
	}, {
		path:     "sub/data_1",
		size:     256,
		seed:     2,
		hash:     "22c8fbc6f57675b9614933fbcbb0f93987385a201004ae4495c6ba6805dbb85d46fcb02222a60d2f151f7346d249027b5fa0684c0ded2e7d0895ece38fce2c6b",
		inDigest: true,
	}}, digest.SHA512)
}

func (s *TestSuiteCalculator) TestExclude() {
	common.Config.Exclude.Absolute = []string{filepath.Join(s.root, "excludeAbs")}
	common.Config.Exclude.Relative = []string{"excludeRel"}
	common.Config.Exclude.Anywhere = []string{"excludeAny"}

	s.verifyDigest([]testCase{{
		path:     "data_1",
		size:     256,
		seed:     1,
		hash:     "352adfeb0dc6e28699635c5911cf33e2e0a86aedf85a5a99bba97749000ae1c7",
		inDigest: true,
	}, {
		path:     "excludeAbs",
		size:     0,
		seed:     1,
		hash:     "",
		inDigest: false,
	}, {
		path:     "excludeRel",
		size:     0,
		seed:     1,
		hash:     "",
		inDigest: false,
	}, {
		path:     "excludeAny",
		size:     0,
		seed:     1,
		hash:     "",
		inDigest: false,
	}, {
		path:     "sub/data_1",
		size:     256,
		seed:     2,
		hash:     "1629705c76a590f2e16b8c42fa0aca9c405401fcfc794399e71f0954f1e0d50e",
		inDigest: true,
	}, {
		path:     "sub/excludeRel",
		size:     256,
		seed:     3,
		hash:     "2e98052dd231a0217464daf09e4a203611b3490845864bf7ea93254c93ca7372",
		inDigest: true,
	}, {
		path:     "sub/excludeAny",
		size:     0,
		seed:     1,
		hash:     "",
		inDigest: false,
	}}, digest.SHA256)
}

func (s *TestSuiteCalculator) verifyDigest(
	testdata []testCase,
	algorithm digest.HashType) {

	createData(s, testdata)

	digest := digest.Calculate(s.root, algorithm)

	count := 0
	for _, dataPoint := range testdata {
		entryHash, entryInDigest := (*digest.Entries)[dataPoint.path]
		assert.Equal(s.T(), dataPoint.inDigest, entryInDigest, dataPoint.path)
		if entryInDigest {
			assert.Equal(s.T(), dataPoint.hash, entryHash, dataPoint.path)
			count++
		}
	}
	require.Equal(s.T(), count, len(*digest.Entries))
}

func createData(
	s *TestSuiteCalculator,
	testdata []testCase) {

	for _, dataPoint := range testdata {
		file := filepath.Join(s.root, dataPoint.path)
		createRandomFile(file, dataPoint.size, dataPoint.seed)
	}
}

func createRandomFile(file string, size, seed int64) {
	err := os.MkdirAll(filepath.Dir(file), 0700)
	if err != nil {
		panic(err)
	}

	out, err := os.Create(file)
	if err != nil {
		panic(err)
	}

	defer out.Close()
	in := rand.New(rand.NewSource(seed))
	io.CopyN(out, in, size)
}
