package metric

import (
	"sync"
)

//type ConcurrentMap struct {
//	IDMetricMap map[string]metric.Metric
//	cmlock      *sync.RWMutex
//}
//
//var containerMetricMap ConcurrentMap
//
//func InitContainerMetricMap() {
//	containerMetricMap = ConcurrentMap{
//		IDMetricMap: make(map[string]metric.Metric),
//		cmlock:      new(sync.RWMutex),
//	}
//}

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
