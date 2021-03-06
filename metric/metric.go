package metric

import (
	"fmt"
	"strings"
	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/golang/glog"
)

func SetGlobalSetting(dclient DockerClient, timeout, force time.Duration) {
	gset = Setting{timeout, force, dclient}
}

func CreateMetric(step time.Duration, fclient Remote, tag string, endpoint string) Metric {
	return Metric{
		Step:     step,
		Client:   fclient,
		Tag:      tag,
		Endpoint: endpoint,
		Stop:     make(chan bool),
	}
}

func (self *Metric) InitMetric(cid string, pid int) (err error) {
	var info map[string]uint64
	if info, err = self.UpdateStats(cid, pid); err == nil {
		self.Last = time.Now()
		self.SaveLast(info)
	} else {
		DeleteContainerMetricMapKey(cid)
	}
	return
}

func (self *Metric) Exit() {
	self.Stop <- true
	close(self.Stop)
}

func (self *Metric) UpdateStats(cid string, pid int) (map[string]uint64, error) {
	info := map[string]uint64{}
	statsChan := make(chan *docker.Stats)
	doneChan := make(chan bool)

	opt := docker.StatsOptions{ID: cid, Stats: statsChan,
		Stream: false, Done: doneChan,
		Timeout: gset.timeout * time.Second}

	go func() {
		if err := gset.dclient.Stats(opt); err != nil {
			if strings.Contains(err.Error(), "No such container") {
				DeleteContainerMetricMapKey(cid)
			} else {
				glog.Warningf("Get stats failed %s: %v", cid[:12], err)
				DeleteContainerMetricMapKey(cid)
			}
		}
	}()

	var stats *docker.Stats = nil
	select {
	case stats = <-statsChan:
		if stats == nil {
			DeleteContainerMetricMapKey(cid)
			return info, fmt.Errorf("Get stats failed: %s", cid[:12])
		}
	case <-time.After(gset.force * time.Second):
		doneChan <- true
		DeleteContainerMetricMapKey(cid)
		return info, fmt.Errorf("Get stats timeout: %s", cid[:12])
	}

	info["docker.cpu.user"] = stats.CPUStats.CPUUsage.UsageInUsermode
	//info["docker.cpu.system"] = stats.CPUStats.CPUUsage.UsageInKernelmode
	info["docker.cpu.system"] = stats.CPUStats.SystemCPUUsage
	info["docker.cpu.usage"] = stats.CPUStats.CPUUsage.TotalUsage
	info["cpuNum"] = uint64(len(stats.CPUStats.CPUUsage.PercpuUsage))
	//FIXME in container it will get all CPUStats
	info["docker.mem.usage"] = stats.MemoryStats.Usage
	info["docker.mem.max_usage"] = stats.MemoryStats.MaxUsage
	info["docker.mem.rss"] = stats.MemoryStats.Stats.Rss
	info["memLimit"] = stats.MemoryStats.Limit

	//FIXME use docker api network data.
	var (
		rxDropped uint64 = 0
		rxBytes   uint64 = 0
		rxErrors  uint64 = 0
		txPackets uint64 = 0
		txDropped uint64 = 0
		rxPackets uint64 = 0
		txErrors  uint64 = 0
		txBytes   uint64 = 0
	)

	for _, v := range stats.Networks {
		rxDropped += v.RxDropped
		rxBytes += v.RxBytes
		rxErrors += v.RxErrors
		txPackets += v.TxPackets
		txDropped += v.TxDropped
		rxPackets += v.RxPackets
		txErrors += v.TxErrors
		txBytes += v.TxBytes
	}

	info["net.rx_dropped"] = uint64(rxDropped)
	info["net.rx_bytes"] = uint64(rxBytes)
	info["net.rx_errors"] = uint64(rxErrors)
	info["net.tx_packets"] = uint64(txPackets)
	info["net.tx_dropped"] = uint64(txDropped)
	info["net.rx_packets"] = uint64(rxPackets)
	info["net.tx_errors"] = uint64(txErrors)
	info["net.tx_bytes"] = uint64(txBytes)

	//FIXME use docker api disk io data.
	var (
		blkRead  uint64 = 0
		blkWrite uint64 = 0
	)
	for _, bioEntry := range stats.BlkioStats.IOServiceBytesRecursive {
		switch strings.ToLower(bioEntry.Op) {
		case "read":
			blkRead = blkRead + bioEntry.Value
		case "write":
			blkWrite = blkWrite + bioEntry.Value
		}
	}
	info["disk.io.read_bytes"] = blkRead
	info["disk.io.write_bytes"] = blkWrite

	return info, nil
}

func (self *Metric) SaveLast(info map[string]uint64) {
	self.Save = map[string]uint64{}
	for k, d := range info {
		self.Save[k] = d
	}
}

func (self *Metric) CalcRate(info map[string]uint64, now time.Time) (rate map[string]float64) {
	rate = map[string]float64{}
	delta := now.Sub(self.Last)
	nano_t := float64(delta.Nanoseconds())
	for k, d := range info {
		switch {
		case strings.HasPrefix(k, "docker.cpu.usage") && d >= self.Save[k]:
			var (
				cpuPercent = 0.0
				// calculate the change for the cpu usage of the container in between readings
				cpuDelta = float64(d - self.Save[k])
				// calculate the change for the entire system between readings
				keysystem   = "docker.cpu.system"
				systemDelta = float64(info[keysystem] - self.Save[keysystem])
			)

			if systemDelta > 0.0 && cpuDelta > 0.0 {
				cpuPercent = (cpuDelta / systemDelta) * float64(info["cpuNum"]) * 100.0
			}
			rate[k] = cpuPercent
		case strings.HasPrefix(k, "docker.cpu.") && d >= self.Save[k]:
			rate[k] = float64(d-self.Save[k]) / nano_t
		case strings.HasPrefix(k, "disk.") && d >= self.Save[k]:
			rate[fmt.Sprintf("docker.%s", k)] = float64(d-self.Save[k]) / nano_t
		case strings.HasPrefix(k, "net.") && d >= self.Save[k]:
			rate[fmt.Sprintf("docker.%s", k)] = float64(d-self.Save[k]) / nano_t
			//rate[fmt.Sprintf("docker.%s", k)] = float64(d)
		case strings.HasPrefix(k, "docker.mem.usage"):
			var memPercent = 0.0
			if info["memLimit"] != 0 {
				memPercent = float64(d) / float64(info["memLimit"]) * 100.0
			}
			rate[k] = float64(memPercent)
		case strings.HasPrefix(k, "docker.mem."):
			rate[k] = float64(d)
		}
	}
	self.Last = now
	return
}

func (self *Metric) Send(rate map[string]float64) error {
	step := int64(self.Step.Seconds())
	timestamp := self.Last.Unix()
	return self.Client.Send(rate, self.Endpoint, self.Tag, timestamp, step)
}
