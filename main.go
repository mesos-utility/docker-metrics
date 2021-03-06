package main

import (
	"flag"
	"strings"
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

	initAndStartWatcher()
}

// init and start watcher
func initAndStartWatcher() {
	metric.InitContainerMetricMap()
	dclient, err := dockerclient.NewDockerClient()
	if err != nil {
		glog.Fatalf("New docker api client error: %v", err)
	}

	metric.SetGlobalSetting(dclient, 6, 8)
	fclient := falcon.CreateFalconClient()
	options := docker.ListContainersOptions{All: true}

	if containers, err := dclient.ListContainers(options); err != nil {
		glog.Errorf("Get container error: %v", err)
	} else {
		for _, container := range containers {
			watcher.AddContainerWatched(dclient, container, fclient)
		}
	}
	var interval int64 = g.Config().Daemon.Interval / 2
	timer := time.NewTicker(time.Duration(interval) * time.Second)

	for {
	REST:
		<-timer.C
		if containers, err := dclient.ListContainers(options); err != nil {
			glog.Errorf("Get container error: %v", err)
			goto REST
		} else {
			for _, container := range containers {
				if _, ok := metric.ContainerMetricMap()[container.ID]; ok {
					if strings.Contains(container.Status, "Up ") {
						continue
					} else {
						metric.DeleteContainerMetricMapKey(container.ID)
					}
				} else {
					if strings.Contains(container.Status, "Up ") {
						if g.Config().Debug {
							glog.Infof("=> Add container: %s", container.ID[:g.IDLEN])
						}
						watcher.AddContainerWatched(dclient, container, fclient)
					} else {
						continue
					}
				}
			}
		}
	}
}
