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

func getCurrentMaven(ctx context.Context) (string, string) {
	const settingsFilePath = "%s/.m2/settings.xml"
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Debug(fmt.Sprintf("error reading home dir : %v", err))
		return "", ""
	}
	userMavenFile := fmt.Sprintf(settingsFilePath, dirname)
	open, err := os.Open(userMavenFile)
	if err != nil {
		log.Debug(fmt.Sprintf("error reading maven settings file : %v", err))
		return "", ""
	}
	doc, err := xmlquery.Parse(open)
	if err != nil {
		panic(err)
	}
	urlNode := xmlquery.FindOne(doc, "//url[1]")
	if urlNode == nil {
		log.Debug(fmt.Sprintf("cannot find a url node on maven config"))
		return "", ""
	}
	parsedUrl, err := url.Parse(urlNode.InnerText())
	if err != nil {
		log.Debug(fmt.Sprintf("error parsing maven repo URL (%s): %v", urlNode.Data, err))
		return "", ""
	}
	serverId := findServerIdByUrl(parsedUrl.String())
	split := strings.Split(parsedUrl.Path, "/")
	if len(split) != 3 {
		return "", ""
	}
	return serverId, split[2]
}
