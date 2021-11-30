package show

import (
	"context"
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/url"
	"os"
	"strings"
)

func getCurrentGolang(ctx context.Context) []repoSelection {
	goProxy := os.Getenv("GOPROXY")
	if goProxy == "" {
		log.Debug("no GOPROXY set")
		return nil
	}
	parsedUrl, err := url.Parse(goProxy)
	if err != nil {
		log.Debug(fmt.Sprintf("error parsing golang repo URL: %v", err))
		return nil
	}
	//remove user so
	parsedUrl.User = nil
	serverId := findServerIdByUrl(parsedUrl.String())
	if serverId == "" {
		return []repoSelection{{
			Unknown: true,
			RepoKey: parsedUrl.String(),
		}}
	}
	split := strings.Split(parsedUrl.Path, "/")
	return []repoSelection{{
		ServerId:    serverId,
		RepoKey:     split[len(split)-1],
		Description: "",
	}}
}
