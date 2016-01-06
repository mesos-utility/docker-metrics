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
)

func handleVersion(displayVersion bool) {
	if displayVersion {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}
}

var cfg = flag.String("c", "cfg.json", "configuration file")
var version = flag.Bool("version", false, "show version")

func main() {
	var transferAddr string
	defer glog.Flush()
	flag.Parse()

	handleVersion(*version)

	// global config
	g.ParseConfig(*cfg)

	metric.InitContainerMetricMap()
	transferAddr = g.Config().Transfer.Addr
	dockerclient, err := dockerclient.NewDockerClient()
	if err != nil {
		glog.Errorf("%v", err)
		return
	}

	metric.SetGlobalSetting(dockerclient, 2, 3, "vnbe", "eth0")
	client := falcon.CreateFalconClient(transferAddr, 5*time.Millisecond)

	if containers, err := dockerclient.ListContainers(docker.ListContainersOptions{All: false}); err != nil {
		glog.Errorf("Get container error: %v", err)
		os.Exit(1)
	} else {
		for _, container := range containers {
			hostname, _ := g.Hostname()
			shortID := container.ID[:g.IDLEN]
			tag := fmt.Sprintf("app=yks-web,id=%s", shortID)
			m := metric.CreateMetric(time.Duration(g.Interval)*time.Second, client, tag, hostname)
			metric.AddContainerMetric(container.ID, m)

			if c, err := dockerclient.InspectContainer(container.ID); err != nil {
				glog.Warningf("%s: %v", container.ID, err)
			} else {
				go watcher(m, c.ID, c.State.Pid)
			}
		}
	}

	for {
	REST:
		time.Sleep(time.Duration(g.Interval/2) * time.Second)
		if containers, err := dockerclient.ListContainers(docker.ListContainersOptions{All: false}); err != nil {
			glog.Errorf("Get container error: %v", err)
			goto REST
		} else {
			for _, container := range containers {
				if _, ok := metric.ContainerMetricMap()[container.ID]; ok {
					continue
				} else {
					fmt.Println("Add ID: ", container.ID)
					hostname, _ := g.Hostname()
					shortID := container.ID[:g.IDLEN]
					tag := fmt.Sprintf("app=yks-web,id=%s", shortID)
					m := metric.CreateMetric(time.Duration(g.Interval)*time.Second, client, tag, hostname)
					metric.AddContainerMetric(container.ID, m)

					if c, err := dockerclient.InspectContainer(container.ID); err != nil {
						glog.Warningf("%s: %v", container.ID, err)
					} else {
						go watcher(m, c.ID, c.State.Pid)
					}
				}
			}
		}
	}
}

func watcher(serv metric.Metric, cid string, pid int) {
	defer serv.Client.Close()
	if err := serv.InitMetric(cid, pid); err != nil {
		glog.Warningf("Fail InitMetric %s: %v", cid, err)
		return
	}

	t := time.NewTicker(serv.Step)
	debug := g.Config().Debug
	defer t.Stop()
	//fmt.Println("begin watch", cid)
	for {
		select {
		case now := <-t.C:
			go func() {
				if info, err := serv.UpdateStats(cid, pid); err == nil {
					if debug {
						glog.Infof("updatestats: %v", cid)
					}
					rate := serv.CalcRate(info, now)
					serv.SaveLast(info)
					// for safe
					//fmt.Println(rate)
					go serv.Send(rate)
				} else {
					glog.Errorf("%v", err)
				}
			}()
		case <-serv.Stop:
			return
		}
	}
}
