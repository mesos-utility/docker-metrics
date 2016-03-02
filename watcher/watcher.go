package watcher

import (
	"fmt"
	"strings"
	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/golang/glog"
	"github.com/mesos-utility/docker-metrics/falcon"
	"github.com/mesos-utility/docker-metrics/g"
	"github.com/mesos-utility/docker-metrics/metric"
)

func AddContainerWatched(dclient *docker.Client, container docker.APIContainers, fclient *falcon.FalconClient) {
	if c, err := dclient.InspectContainer(container.ID); err != nil {
		glog.Warningf("%s: %v", container.ID, err)
	} else {
		hostname, _ := g.Hostname()
		// format tags
		shortID := container.ID[:g.IDLEN]
		tag := getTagFromContainer(c)
		tags := fmt.Sprintf("%s,id=%s", tag, shortID)
		attachTags := strings.TrimSpace(g.Config().AttachTags)
		if attachTags != "" {
			tags += fmt.Sprintf(",%s", attachTags)
		}
		tags = strings.Trim(tags, ",")
		// get interval
		interval := g.Config().Daemon.Interval
		m := metric.CreateMetric(time.Duration(interval)*time.Second, fclient, tags, hostname)
		metric.AddContainerMetric(container.ID, m)
		go watcher(m, c.ID, c.State.Pid)
	}
}

// Get tag from container: app=MARATHON_APP_ID or image=Imagename
func getTagFromContainer(ct *docker.Container) (tag string) {
	tag = ""
	for _, v := range ct.Config.Env {
		if strings.HasPrefix(v, "MARATHON_APP_ID") {
			tag = strings.Split(v, "=")[1]
			if strings.HasPrefix(tag, "/") {
				tag = tag[1:]
			}
			tag = fmt.Sprintf("app=%s", tag)
			break
		}
	}

	if tag == "" {
		imglen = len(ct.Config.Image)

		if strings.Contains(ct.Config.Image, "/") {
			tmparray := strings.Split(ct.Config.Image, "/")
			tag = tmparray[len(tmparray)-1]
		} else if imglen > 12 {
			tag = ct.Config.Image[:12]
		} else {
			tag = ct.Config.Image
		}

		if tag != "" {
			tag = fmt.Sprintf("image=%s", tag)
		}
	}

	return tag
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
