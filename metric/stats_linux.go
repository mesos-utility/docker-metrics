// +build linux

package metric

// Get network stats from /proc/{pid}/net/dev
func (self *Metric) getNetStats(result map[string]uint64) (err error) {
	return
}

// Get disk stats from /proc/{pid}/io
// may failed by permission denied.
func (self *Metric) getDiskStats(result map[string]uint64) (err error) {
	return
}
