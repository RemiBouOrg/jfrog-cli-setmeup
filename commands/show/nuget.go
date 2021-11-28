package show

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/url"
	"os/exec"
	"strings"
)

func getCurrentNuget(ctx context.Context) []repoSelection {
	command := exec.Command("nuget", "sources", "List",
		"-Format", "Short",
	)
	bufferString := bytes.NewBufferString("")
	command.Stdout = bufferString
	err := command.Run()
	if err != nil {
		log.Debug(fmt.Sprintf("error running nuget source List %v", err))
		return nil
	}
	res := make([]repoSelection, 0)
	scanner := bufio.NewScanner(bufferString)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		if line[0] != "E" {
			continue
		}
		urlRepository := line[1]
		serverId := findServerIdByUrl(urlRepository)
		if serverId == "" {
			res = append(res, repoSelection{
				unknown: true,
				repoKey: urlRepository,
			})
			continue

		}
		parsedUrl, err := url.Parse(urlRepository)
		if err != nil {
			log.Debug(fmt.Sprintf("error parsing nuget repo URL (%s): %v", urlRepository, err))
			return nil
		}
		split := strings.Split(parsedUrl.Path, "/")
		res = append(res, repoSelection{
			serverId:    serverId,
			repoKey:     split[len(split)-1],
			description: "",
		})
	}
	return res
}
