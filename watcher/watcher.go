package watcher

import (
	"fmt"
	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/golang/glog"
	"github.com/mesos-utility/docker-metrics/falcon"
	"github.com/mesos-utility/docker-metrics/g"
	"github.com/mesos-utility/docker-metrics/metric"
)

func AddContainerWatched(dclient *docker.Client, container docker.APIContainers, client *falcon.FalconClient) {
	hostname, _ := g.Hostname()
	shortID := container.ID[:g.IDLEN]
	tag := fmt.Sprintf("app=yks-web,id=%s", shortID)
	m := metric.CreateMetric(time.Duration(g.Interval)*time.Second, client, tag, hostname)
	metric.AddContainerMetric(container.ID, m)

	if c, err := dclient.InspectContainer(container.ID); err != nil {
		glog.Warningf("%s: %v", container.ID, err)
	} else {
		go watcher(m, c.ID, c.State.Pid)
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
