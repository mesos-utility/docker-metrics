package dockerclient

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fsouza/go-dockerclient"
	"github.com/mesos-utility/docker-metrics/g"
)

var certs = []string{"cert.pem", "key.pem", "ca.pem"}

// new docker client use go-dockerclient, ENV TLS and Common.
func NewDockerClient() (client *docker.Client, err error) {
	daemonAddr := g.Config().Daemon.Addr
	certDir := g.Config().Daemon.CertDir

	if daemonAddr == "" {
		client, err = docker.NewClientFromEnv()
	} else {
		if !strings.HasPrefix(daemonAddr, "tcp://") {
			return nil, fmt.Errorf("Please check docker addr in cfg.json!!!")
		}

		if _, err = g.CheckFilesExist(certDir, certs); err == nil {
			cert := filepath.Join(certDir, "cert.pem")
			key := filepath.Join(certDir, "key.pem")
			ca := filepath.Join(certDir, "ca.pem")
			client, err = docker.NewTLSClient(daemonAddr, cert, key, ca)
		} else {
			client, err = docker.NewClient(daemonAddr)
		}
	}

	return client, err
}
