package driver

import (
	build "github.com/Azure/acr-builder/pkg"
	"github.com/Azure/acr-builder/pkg/commands"
	"github.com/Azure/acr-builder/pkg/constants"
)

func getSource(workingDir,
	gitURL, gitBranch, gitHeadRev, gitXToken, gitPATokenUser,
	gitPAToken string) (source build.Source, err error) {

	var gitFactory, localFactory, selected *factory

	gitFactory, err = newFactory(constants.SourceNameGit,
		func() (build.Source, error) {
			var gitCred commands.GitCredential
			if gitXToken != "" {
				gitCred = commands.NewGitXToken(gitXToken)
			} else if gitPATokenUser != "" {
				var err error
				gitCred, err = commands.NewGitPersonalAccessToken(gitPATokenUser, gitPAToken)
				if err != nil {
					return nil, err
				}
			}
			return commands.NewGitSource(gitURL, gitBranch, gitHeadRev, workingDir, gitCred), nil
		},
		[]parameter{
			{name: constants.ArgNameGitURL, value: gitURL},
		},
		[]parameter{
			{name: constants.ArgNameGitBranch, value: gitBranch},
			{name: constants.ArgNameGitHeadRev, value: gitHeadRev},
			{name: constants.ArgNameGitXToken, value: gitXToken},
			{name: constants.ArgNameGitPATokenUser, value: gitPATokenUser},
			{name: constants.ArgNameGitPAToken, value: gitPAToken},
		},
	)
	if err != nil {
		return
	}

	localFactory, err = newFactory(constants.SourceNameLocal,
		func() (build.Source, error) {
			return commands.NewLocalSource(workingDir), nil
		}, nil, nil)
	if err != nil {
		return
	}

	selected, err = decide("sources", localFactory, gitFactory)
	if err != nil {
		return
	}

	return selected.create()
}
