package metric

import (
	"sync"

	"github.com/golang/glog"
	"github.com/mesos-utility/docker-metrics/g"
)

var (
	idMetricMap map[string]Metric
	cmlock      = new(sync.RWMutex)
)

func InitContainerMetricMap() {
	cmlock.Lock()
	defer cmlock.Unlock()

	idMetricMap = make(map[string]Metric)
}

func ContainerMetricMap() map[string]Metric {
	cmlock.RLock()
	defer cmlock.RUnlock()

	return idMetricMap
}

func DeleteContainerMetricMapKey(key string) bool {
	if _, ok := idMetricMap[key]; !ok {
		return false
	} else {
		cmlock.Lock()
		defer cmlock.Unlock()

		if g.Config().Debug {
			glog.Infof("<= Delete container: %s", key[:g.IDLEN])
		}
		delete(idMetricMap, key)
	}

	return true
}

func AddContainerMetric(key string, value Metric) bool {
	if _, ok := idMetricMap[key]; ok {
		return false
	} else {
		cmlock.Lock()
		defer cmlock.Unlock()

		idMetricMap[key] = value
	}

	return true
}
