package jfrogconfig

import (
	"errors"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/commons"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const jfrogConfFilePath = "./.jfrog-setmeup.yaml"

func FindRepoKeys(c *components.Context, serverDetails *config.ServerDetails) ([]string, error) {
	if len(c.Arguments) >= 1 {
		repoArg := c.Arguments[0]
		if len(repoArg) > 0 {
			return []string{repoArg}, nil
		}
	}

	repo, err := commons.FindRepo(serverDetails)
	if err != nil {
		return nil, err
	}

	if repo != nil {
		return []string{repo.Key}, nil
	}

	return nil, errors.New("wrong number of arguments. Expected repository key or use jfrog setmeup init <repo-key> to store in the source control")
}

func extractValues(repoKeys *RepoTypeToName) []string {
	allValues := make([]string, 0, len(*repoKeys))
	for _, value := range *repoKeys {
		allValues = append(allValues, value)
	}
	return allValues
}

func findRepoKeyFromConfFile(envName string) (*RepoTypeToName, error) {
	confFile, err := ReadCurrentConfFile()
	if err != nil {
		return nil, err
	}

	repoTypeToName := confFile[envName]
	return &repoTypeToName, nil
}

func ReadCurrentConfFile() (JFrogConfFile, error) {
	confFileContent, err := fileutils.ReadFile(jfrogConfFilePath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	confFile := JFrogConfFile{}
	err = yaml.Unmarshal(confFileContent, &confFile)
	if err != nil {
		return nil, os.Remove(jfrogConfFilePath)
	}
	return confFile, nil
}

func WriteConfigFile(confContent []byte) error {
	return ioutil.WriteFile(jfrogConfFilePath, confContent, 0644)
}
