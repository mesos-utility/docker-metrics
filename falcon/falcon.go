package falcon

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"net/rpc"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/mesos-utility/docker-metrics/g"
	"github.com/open-falcon/common/model"
	"github.com/toolkits/net"
)

func CreateFalconClient() *FalconClient {
	transferAddr := g.Config().Transfer.Addr
	timeout := time.Duration(g.Config().Transfer.Timeout) * time.Millisecond
	transfer := true

	if strings.HasPrefix(transferAddr, "http") {
		transfer = false
	}

	return &FalconClient{
		transfer:  transfer,
		RpcServer: transferAddr,
		Timeout:   timeout,
	}
}

type FalconClient struct {
	sync.Mutex
	rpcClient *rpc.Client
	transfer  bool
	RpcServer string
	Timeout   time.Duration
}

func (self *FalconClient) Close() error {
	if self.rpcClient != nil {
		self.rpcClient.Close()
		self.rpcClient = nil
	}
	return nil
}

func (self *FalconClient) insureConn() error {
	if self.rpcClient != nil {
		return nil
	}

	var err error
	var retry int = 1

	for {
		if self.rpcClient != nil {
			return nil
		}

		self.rpcClient, err = net.JsonRpcClient("tcp", self.RpcServer, self.Timeout)
		if err == nil {
			return nil
		}

		glog.Warningf("Metrics rpc dial fail %s", err)
		if retry > 5 {
			return err
		}

		time.Sleep(time.Duration(math.Pow(2.0, float64(retry))) * time.Second)
		retry++
	}
}

func (self *FalconClient) call(method string, args interface{}, reply interface{}) error {

	self.Lock()
	defer self.Unlock()

	if err := self.insureConn(); err != nil {
		return err
	}

	timeout := time.Duration(50 * time.Second)
	done := make(chan error)

	go func() {
		err := self.rpcClient.Call(method, args, reply)
		done <- err
	}()

	select {
	case <-time.After(timeout):
		glog.Warningf("Metrics rpc call timeout %s %s", self.rpcClient, self.RpcServer)
		self.Close()
	case err := <-done:
		if err != nil {
			self.Close()
			return err
		}
	}
	return nil
}

func (self *FalconClient) Send(data map[string]float64, endpoint, tag string, timestamp, step int64) error {
	metrics := []*model.MetricValue{}
	var metric *model.MetricValue
	for k, v := range data {
		metric = &model.MetricValue{
			Endpoint:  endpoint,
			Metric:    k,
			Value:     v,
			Step:      step,
			Type:      "GAUGE",
			Tags:      tag,
			Timestamp: timestamp,
		}
		metrics = append(metrics, metric)
	}
	//log.Printf("%v\n", metrics)

	if len(metrics) == 0 {
		return nil
	}

	var resp model.TransferResponse
	if g.Config().Transfer.Enable {
		if self.transfer {
			if err := self.call("Transfer.Update", metrics, &resp); err != nil {
				return err
			}
		} else {
			if err := PostToAgent(metrics); err != nil {
				return err
			}
		}
	} else {
		glog.Infoln("=> <Total=%d> %v\n", len(metrics), metrics[0])
	}

	if g.Config().Debug {
		glog.Infof("%s: %v %v", endpoint, timestamp, &resp)
	}
	return nil
}

//      PostToAgent     ->  http://127.0.0.1:1988/v1/push
//      Send            ->  127.0.0.1:8433
func PostToAgent(metrics []*model.MetricValue) error {
	if len(metrics) == 0 {
		return nil
	}

	debug := g.Config().Debug

	if debug {
		glog.Infof("=> <Total=%d> %v\n", len(metrics), metrics[0])
	}

	contentJson, err := json.Marshal(metrics)
	if err != nil {
		glog.Warningf("Error for PostToAgent json Marshal: %v", err)
		return err
	}
	contentReader := bytes.NewReader(contentJson)
	req, err := http.NewRequest("POST", g.Config().Transfer.Addr, contentReader)
	if err != nil {
		glog.Warningf("Error for PostToAgent in NewRequest: %v", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		glog.Warningf("Error for PostToAgent in http client Do: %v", err)
		return err
	}
	defer resp.Body.Close()

	if debug {
		glog.Infof("<= %v", resp.Body)
	}

	return nil
}
