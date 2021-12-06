package repository

import (
	"bytes"
	"context"
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/pkg/errors"
	"os/exec"
	"strings"
)

func handleGo(ctx context.Context, configuration SetMeUpConfiguration) error {
	urlParts := strings.SplitN(configuration.ServerDetails.ArtifactoryUrl, "://", 2)
	if len(urlParts) < 2 {
		return fmt.Errorf("cannot parse artifactory url %s", configuration.ServerDetails.ArtifactoryUrl)
	}

	artiUrl := fmt.Sprintf("%s://%s:%s@%s", urlParts[0], configuration.ServerDetails.User, configuration.ServerDetails.Password, urlParts[1])

	goProxy := fmt.Sprintf("%s/api/go/%s", artiUrl , configuration.RepoDetails.Key)

	goBinPath, err := exec.LookPath("go")
	if err != nil {
		return fmt.Errorf("failed to detect go executable with error %w", err)
	}
	if err = setGoProxy(ctx, goBinPath, goProxy); err != nil {
		return err
	}

	log.Info(fmt.Sprintf("Repo %s set as GOPROXY", configuration.RepoDetails.Key))
	return nil
}

func setGoProxy(ctx context.Context, goBinPath, goProxy string) error {
	command := exec.CommandContext(ctx, goBinPath, "env", "-w", fmt.Sprintf("GOPROXY=%s", goProxy))
	bufferString := bytes.NewBufferString("")
	command.Stderr = bufferString
	err := command.Run()
	stderr := bufferString.String()
	if err != nil {
		return errors.Wrap(err, stderr)
	}

	if len(stderr) > 0 {
		return errors.New(stderr)
	}

	return nil
}
