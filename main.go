package main

import (
	"flag"
	"fmt"
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
	defer glog.Flush()
	flag.Parse()

	g.HandleVersion(*version)

	// global config
	g.ParseConfig(*cfg)

	metric.InitContainerMetricMap()
	dclient, err := dockerclient.NewDockerClient()
	if err != nil {
		glog.Fatalf("New docker api client error: %v", err)
	}

	metric.SetGlobalSetting(dclient, 2, 3, "vnbe", "eth0")
	fclient := falcon.CreateFalconClient()

	if containers, err := dclient.ListContainers(docker.ListContainersOptions{All: false}); err != nil {
		glog.Errorf("Get container error: %v", err)
	} else {
		for _, container := range containers {
			watcher.AddContainerWatched(dclient, container, fclient)
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
					watcher.AddContainerWatched(dclient, container, fclient)
				}
			}
		}
	}
}
