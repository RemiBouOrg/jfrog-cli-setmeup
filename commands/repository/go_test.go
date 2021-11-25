package repository

import (
	"bytes"
	"context"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func Test_handleGo(t *testing.T) {
	type args struct {
		configuration SetMeUpConfiguration
	}
	tests := []struct {
	name    string
	args    args
	wantErr bool
	want    string
}{
		{
			name: "should have go proxy for valid parameters - https",
			args: args{
				configuration: SetMeUpConfiguration{
					serverDetails: &config.ServerDetails{ArtifactoryUrl: "https://example.com/artifactory", User: "foo", Password: "bar"},
					repoDetails:   &RepoDetails{Key: "go-local"},
				},
			},
			want: "https://foo:bar@example.com/artifactory/api/go/go-local",
			wantErr: false,
		},
		{
			name: "should have go proxy for valid parameters - http",
			args: args{
				configuration: SetMeUpConfiguration{
					serverDetails: &config.ServerDetails{ArtifactoryUrl: "http://example.com/artifactory", User: "foo", Password: "bar"},
					repoDetails:   &RepoDetails{Key: "go-local"},
				},
			},
			want: "http://foo:bar@example.com/artifactory/api/go/go-local",
			wantErr: false,
		},
		{
			name: "should fail on invalid arti url",
			args: args{
				configuration: SetMeUpConfiguration{
					serverDetails: &config.ServerDetails{ArtifactoryUrl: "example.com/artifactory", User: "foo", Password: "bar"},
					repoDetails:   &RepoDetails{Key: "go-local"},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		_ = os.Unsetenv("GOPROXY")
		t.Run(tt.name, func(t *testing.T) {
			err := handleGo(context.Background(), tt.args.configuration);
			if (err != nil) != tt.wantErr {
				t.Errorf("handleGo() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil {
				command := exec.Command(os.Getenv("GOROOT")+"/bin/go", "env", "GOPROXY")
				bufferString := bytes.NewBufferString("")
				command.Stdout = bufferString

				err := command.Run()
				require.NoError(t, err)
				assert.Equal(t, tt.want, strings.Trim(bufferString.String(), "\"\n"))
			}
		})
	}
}
