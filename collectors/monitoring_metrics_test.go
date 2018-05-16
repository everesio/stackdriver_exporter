package collectors

import (
	"testing"
	"google.golang.org/api/monitoring/v3"
	"github.com/prometheus/client_golang/prometheus"
	"fmt"
)

func TestFillHistogramMetricsLabels(t *testing.T) {
	metricDescriptor := &monitoring.MetricDescriptor{}

	ch := make(chan prometheus.Metric, 256)

	timeSeriesMetrics := &TimeSeriesMetrics{
		metricDescriptor:  metricDescriptor,
		ch:                ch,
		fillMissingLabels: true,
		constMetrics:      make(map[string][]ConstMetric),
		histogramMetrics:  make(map[string][]HistogramMetric),
	}
	timeSeries := &monitoring.TimeSeries{Resource: &monitoring.MonitoredResource{Type: "https_lb_rule_loadbalancing_googleapis_com"}, Metric: &monitoring.Metric{Type: "https_total_latencies"}}
	timeSeriesMetrics.CollectNewConstHistogram(timeSeries, []string {"backend_name", "backend_scope"},  &monitoring.Distribution{}, nil, []string {"k8s-ig--635e72c38338af7e", "europe-west1-c"})
	timeSeriesMetrics.CollectNewConstHistogram(timeSeries, []string {"backend_scope"},  &monitoring.Distribution{}, nil, []string {"europe-west1-c"})
	timeSeriesMetrics.completeHistogramMetrics()

	for _, vs := range timeSeriesMetrics.histogramMetrics {
		if len(vs) > 1 {
			var needFill bool
			for i := 1; i < len(vs); i++ {
				if vs[0].keysHash != vs[i].keysHash {
					needFill = true
				}
			}
			if needFill {
				vs = fillHistogramMetricsLabels(vs)
			}
		}
		for _, v := range vs {
			fmt.Println(v.labelKeys, "-", len(v.labelKeys), "#", v.labelValues, "-", len(v.labelValues))
		}
	}
}