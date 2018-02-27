package driver

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	build "github.com/Azure/acr-builder/pkg"
	"github.com/Azure/acr-builder/pkg/commands"
)

// happy cases are tested in build_test.go
// we will mainly test negative case

type getSourceTestCase struct {
	workingDir     string
	gitURL         string
	gitBranch      string
	gitHeadRev     string
	gitXToken      string
	gitPATokenUser string
	gitPAToken     string
	expectedError  string
	expectedSource build.Source
}

func TestGetSourceLocal(t *testing.T) {
	target := "home"
	testGetSource(t, getSourceTestCase{
		workingDir:     target,
		expectedSource: commands.NewLocalSource(target),
	})
}

func testGetSource(t *testing.T, tc getSourceTestCase) {
	source, err := getSource(tc.workingDir,
		tc.gitURL, tc.gitBranch, tc.gitHeadRev, tc.gitXToken, tc.gitPATokenUser, tc.gitPAToken)

	if tc.expectedError != "" {
		assert.NotNil(t, err)
		assert.Regexp(t, regexp.MustCompile(tc.expectedError), err.Error())
		return
	}

	assert.Nil(t, err)
	assert.Equal(t, tc.expectedSource, source)
}
