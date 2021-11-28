package show

import (
	"context"
	"fmt"
	"github.com/antchfx/xmlquery"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/url"
	"os"
	"strings"
)

func getCurrentMaven(ctx context.Context) []repoSelection {
	const settingsFilePath = "%s/.m2/settings.xml"
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Debug(fmt.Sprintf("error reading home dir : %v", err))
		return nil
	}
	userMavenFile := fmt.Sprintf(settingsFilePath, dirname)
	open, err := os.Open(userMavenFile)
	if err != nil {
		log.Debug(fmt.Sprintf("error reading maven settings file : %v", err))
		return nil
	}
	doc, err := xmlquery.Parse(open)
	if err != nil {
		panic(err)
	}
	repositoryNodes := xmlquery.Find(doc, "//repository")
	if repositoryNodes == nil {
		log.Debug(fmt.Sprintf("cannot find a repositoryNodes node on maven config"))
		return nil
	}
	res := make([]repoSelection, 0)
	for _, repositoryNode := range repositoryNodes {
		urlRepository := repositoryNode.SelectElement("//url").InnerText()
		repositoryId := repositoryNode.SelectElement("//id").InnerText()
		serverId := findServerIdByUrl(urlRepository)
		if serverId == "" {
			res = append(res, repoSelection{
				unknown:     true,
				serverId:    urlRepository,
				description: repositoryId,
			})
			continue
		}
		parsedUrl, err := url.Parse(urlRepository)
		if err != nil {
			log.Debug(fmt.Sprintf("error parsing maven repo URL (%s): %v", urlRepository, err))
			return nil
		}
		split := strings.Split(parsedUrl.Path, "/")
		if len(split) != 3 {
			res = append(res, repoSelection{
				unknown:     true,
				serverId:    urlRepository,
				description: repositoryId,
			})
			continue
		}
		res = append(res, repoSelection{
			serverId:    serverId,
			repoKey:     split[2],
			description: repositoryId,
		})
	}
	return res
}
