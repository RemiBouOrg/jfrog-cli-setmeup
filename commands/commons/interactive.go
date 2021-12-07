package commons

import (
	"fmt"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-cli-plugin-template/commands/artifactory"
	"github.com/manifoldco/promptui"
	"strings"
)

type FindRepoService interface {
	FindRepo(serverDetails *config.ServerDetails) (*artifactory.RepoDetails, error)
}

type findRepoService struct {
}

func NewFindRepoService() *findRepoService {
	return &findRepoService{}
}

func (frs *findRepoService) FindRepo(serverDetails *config.ServerDetails) (*artifactory.RepoDetails, error) {
	repos, err := artifactory.GetAllRepoNames(serverDetails)
	if err != nil {
		return nil, fmt.Errorf("unable to get all repos: %w", err)
	}

	return frs.runInteractiveMenu("", "Repository key", repos)
}

func unzipRepo(details []*artifactory.RepoDetails) []string {
	display := make([]string, 0, len(details))
	for _, detail := range details {
		display = append(display, fmt.Sprintf("\u001b[31m%s\u001b[37m %s", detail.PackageType, detail.Key))
	}

	return display
}

func (*findRepoService) runInteractiveMenu(selectionHeader string, selectionLabel string, values []*artifactory.RepoDetails) (*artifactory.RepoDetails, error) {
	if selectionHeader != "" {
		fmt.Println(selectionHeader)
	}
	selectMenu := promptui.Select{
		Label: selectionLabel,
		Items: values,
		Searcher: func(input string, index int) bool {
			curr := values[index]
			return strings.Index(curr.Key+" "+curr.PackageType, input) >= 0
		},
		StartInSearchMode: true,
		Templates: &promptui.SelectTemplates{
			Active:   "\U0001F438 \u001b[4m{{ .PackageType | red }} :: \u001b[4m{{ .Key | yellow }}",
			Inactive: "{{ .PackageType | red }} :: {{ .Key | yellow }}",
			Selected: "{{ .PackageType | red }} :: {{ .Key | yellow }}",
		},
	}
	selected, _, err := selectMenu.Run()
	return values[selected], err
}
