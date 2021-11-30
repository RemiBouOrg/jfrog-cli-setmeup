package show

import (
	"bytes"
	"context"
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/url"
	"os/exec"
	"strings"
)

func getCurrentNpm(ctx context.Context) []repoSelection {
	command := exec.Command("npm", "config", "get", "registry")
	bufferString := bytes.NewBufferString("")
	command.Stdout = bufferString
	err := command.Run()
	if err != nil {
		log.Debug(fmt.Sprintf("error occured when getting npm registry %v", err))
		return nil
	}
	urlRepository := bufferString.String()
	serverId := findServerIdByUrl(urlRepository)
	if serverId == "" {
		return []repoSelection{
			{
				Unknown: true,
				RepoKey: urlRepository,
			},
		}
	}
	parsedUrl, err := url.Parse(strings.TrimSuffix(urlRepository, "\n"))
	if err != nil {
		log.Debug(fmt.Sprintf("error parsing npm repo URL (%s): %v", urlRepository, err))
		return nil
	}
	split := strings.Split(parsedUrl.Path, "/")
	return []repoSelection{
		{
			ServerId:    serverId,
			RepoKey:     split[len(split)-1],
			Description: "",
		},
	}
}
