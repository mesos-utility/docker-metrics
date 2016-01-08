package dockerclient

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fsouza/go-dockerclient"
	"github.com/mesos-utility/docker-metrics/g"
)

var certs = []string{"cert.pem", "key.pem", "ca.pem"}

// new docker client use go-dockerclient, ENV TLS and Common.
func NewDockerClient() (client *docker.Client, err error) {
	dockerAddr := g.Config().Daemon.Addr
	certDir := g.Config().Daemon.CertDir

	if dockerAddr == "" {
		client, err = docker.NewClientFromEnv()
	} else {
		if !strings.HasPrefix(dockerAddr, "tcp://") {
			return nil, errors.New("Please check docker addr in cfg.json!!!")
		}

		if _, err := g.CheckFilesExist(certDir, certs); err == nil {
			cert := fmt.Sprintf("%s/cert.pem", certDir)
			key := fmt.Sprintf("%s/key.pem", certDir)
			ca := fmt.Sprintf("%s/ca.pem", certDir)
			client, err = docker.NewTLSClient(dockerAddr, cert, key, ca)
		} else {
			client, err = docker.NewClient(dockerAddr)
		}
	}

	return client, err
}
