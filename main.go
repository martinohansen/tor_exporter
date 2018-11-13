package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/martinohansen/tor_exporter/collector"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/yawning/bulb"
)

var (
	listenAddress = flag.String(
		"web.listen-address", ":9105",
		"Address to listen on for metrics scraping",
	)
	torControlPort = flag.String(
		"tor.control-port", "",
		"Address and port of the tor daemon",
	)
	torControlSocket = flag.String(
		"tor.control-socket", "",
		"path to the control socket of the tor daemon",
	)
	torPassword = flag.String(
		"tor.password", "",
		"password to authenticate to the tor port",
	)
)

func main() {
	flag.Parse()

	err := realMain()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func realMain() error {
	if *torControlPort == "" && *torControlSocket == "" {
		return errors.New("either --tor.control-port or --tor.control-socket must be passed")
	}

	// Make a connection just to test
	c, err := connect(*torControlSocket, *torControlPort, *torPassword)
	if err != nil {
		return err
	}
	c.Close()

	exporter := Exporter{
		cSocket:  *torControlSocket,
		cPort:    *torControlPort,
		password: *torPassword,
	}
	prometheus.MustRegister(&exporter)

	http.Handle("/metrics", prometheus.Handler())

	return http.ListenAndServe(*listenAddress, nil)
}

type Exporter struct {
	cSocket  string
	cPort    string
	password string
	isError  prometheus.Gauge
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.NumCircuitsDesc
	ch <- collector.NumStreamsDesc
	ch <- collector.NumOrconnsDesc
	ch <- collector.TrafficReadDesc
	ch <- collector.TrafficWrittenDesc
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	c, err := connect(e.cSocket, e.cPort, e.password)
	if err != nil {
		e.isError.Set(1)
		return
	}

	defer c.Close()

	if err := collector.ScrapeTraffic(c, ch); err != nil {
		log.Errorln(err)
	}

	if err := collector.ScrapeCircuits(c, ch); err != nil {
		log.Errorln(err)
	}
}

func connect(cSocket, cPort, password string) (*bulb.Conn, error) {
	var c *bulb.Conn
	var err error

	if cSocket != "" {
		c, err = bulb.Dial("unix", cSocket)
	} else if cPort != "" {
		c, err = bulb.Dial("tcp4", cPort)
	} else {
		return nil, errors.New("no endpoint specified")
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to tor")
	}

	//c.Debug(true)
	if err = c.Authenticate(""); err != nil {
		return nil, errors.Wrap(err, "Could not authenticate")
	}

	return c, nil
}
