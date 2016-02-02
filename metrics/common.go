package metrics

import (
	"fmt"
	"os"
	"time"

	"gop.kg/src-d/domain@v3/container"

	"github.com/prometheus/client_golang/prometheus"
)

const subsystem = "rovers"

func init() {
	go func() {
		for _ = range time.Tick(container.Config.Prometheus.PushFrequency) {
			Push()
		}
	}()
}

func Push() {
	host, _ := os.Hostname()
	err := prometheus.Push(subsystem, host, container.Config.Prometheus.PushGatewayUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Error pushing collectors to prometheus pushgateway: %s\n", err,
		)
	}
}
