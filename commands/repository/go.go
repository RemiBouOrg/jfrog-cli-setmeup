package repository

import (
	"bytes"
	"context"
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/pkg/errors"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"strings"
)

func handleGo(ctx context.Context, configuration SetMeUpConfiguration) error {
	urlParts := strings.SplitN(configuration.ServerDetails.ArtifactoryUrl, "://", 2)
	if len(urlParts) < 2 {
		return fmt.Errorf("cannot parse artifactory url %s", configuration.ServerDetails.ArtifactoryUrl)
	}

	artiUrl := fmt.Sprintf("%s://%s:%s@%s", urlParts[0], configuration.ServerDetails.User, configuration.ServerDetails.Password, urlParts[1])

	goProxy := fmt.Sprintf("%s/api/go/%s", artiUrl , configuration.RepoDetails.Key)

	paths := []string{path.Join(os.Getenv("GOROOT"), "bin"), "", "/usr/local/bin/", "/usr/bin/"}

	var err error
	for _, p := range paths {
		err = setGoProxy(ctx, p, goProxy)
		if err == nil {
			break
		}

		switch errors.Cause(err).(type) {
		case *fs.PathError:
			continue
		default:
			return err
		}
	}

	if err != nil {
		return err
	}

	log.Info(fmt.Sprintf("Repo %s set as GOPROXY", configuration.RepoDetails.Key))
	return nil
}

func setGoProxy(ctx context.Context, pathPrefix, goProxy string) error {
	if !strings.HasSuffix(pathPrefix, "/") {
		pathPrefix += "/"
	}

	command := exec.CommandContext(ctx, fmt.Sprintf("%sgo", pathPrefix), "env", "-w", fmt.Sprintf("GOPROXY=\"%s\"", goProxy))
	bufferString := bytes.NewBufferString("")
	command.Stderr = bufferString
	err := command.Run()
	if err != nil {
		return errors.Wrap(err, bufferString.String())
	}

	return nil
}
