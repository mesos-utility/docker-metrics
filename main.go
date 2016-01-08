package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/golang/glog"
	"github.com/mesos-utility/docker-metrics/dockerclient"
	"github.com/mesos-utility/docker-metrics/falcon"
	"github.com/mesos-utility/docker-metrics/g"
	"github.com/mesos-utility/docker-metrics/metric"
	"github.com/mesos-utility/docker-metrics/watcher"
)

var cfg = flag.String("c", "cfg.json", "configuration file")
var version = flag.Bool("version", false, "show version")

func main() {
	var transferAddr string
	defer glog.Flush()
	flag.Parse()

	g.HandleVersion(*version)

	// global config
	g.ParseConfig(*cfg)

	metric.InitContainerMetricMap()
	transferAddr = g.Config().Transfer.Addr
	dclient, err := dockerclient.NewDockerClient()
	if err != nil {
		glog.Errorf("%v", err)
		return
	}

	metric.SetGlobalSetting(dclient, 2, 3, "vnbe", "eth0")
	client := falcon.CreateFalconClient(transferAddr, 5*time.Millisecond)

	if containers, err := dclient.ListContainers(docker.ListContainersOptions{All: false}); err != nil {
		glog.Errorf("Get container error: %v", err)
		os.Exit(1)
	} else {
		for _, container := range containers {
			watcher.AddContainerWatched(dclient, container, client)
		}
	}

	for {
	REST:
		interval := g.Config().Daemon.Interval / 2
		time.Sleep(time.Duration(interval) * time.Second)
		if containers, err := dclient.ListContainers(docker.ListContainersOptions{All: false}); err != nil {
			glog.Errorf("Get container error: %v", err)
			goto REST
		} else {
			for _, container := range containers {
				if _, ok := metric.ContainerMetricMap()[container.ID]; ok {
					continue
				} else {
					fmt.Println("Add ID: ", container.ID)
					watcher.AddContainerWatched(dclient, container, client)
				}
			}
		}
	}
}
