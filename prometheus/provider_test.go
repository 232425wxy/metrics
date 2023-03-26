package prometheus

import (
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/232425wxy/metrics"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/require"
)

func TestProvider(t *testing.T) {
	p := &Provider{}

	registry := prom.NewRegistry()
	prom.DefaultRegisterer = registry
	prom.DefaultGatherer = registry
	server := httptest.NewServer(promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	client := server.Client()

	counterOpts := metrics.CounterOpts{
		Namespace:  "counter_namespace",
		Subsystem:  "counter_subsystem",
		Name:       "counter_name",
		Help:       "counter_help",
		LabelNames: []string{"p2p", "consensus"},
	}

	counter := p.NewCounter(counterOpts)
	counter.With("p2p", "gossip", "consensus", "pos").Add(0.1)
	counter.With("p2p", "gossip", "consensus", "pow").Add(0.8)

	resp, err := client.Get(fmt.Sprintf("http://%s/metrics", server.Listener.Addr().String()))
	require.NoError(t, err)
	defer resp.Body.Close()

	bz, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	t.Log(string(bz))

	gaugeOpts := metrics.GaugeOpts{
		Namespace:  "gauge_namespace",
		Subsystem:  "gauge_subsystem",
		Name:       "gauge_name",
		Help:       "gauge_help",
		LabelNames: []string{"cpu", "gpu"},
	}

	gauge := p.NewGauge(gaugeOpts)
	gauge.With("cpu", "intel", "gpu", "3090").Set(99.9)
	gauge.With("cpu", "intel", "gpu", "4090").Add(88.9)
	gauge.With("cpu", "intel", "gpu", "3090").Add(99.9)

	resp, err = client.Get(fmt.Sprintf("http://%s/metrics", server.Listener.Addr().String()))
	require.NoError(t, err)
	defer resp.Body.Close()

	bz, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	t.Log(string(bz))

	histogramOpts := metrics.HistogramOpts{
		Namespace:  "histogram_namespace",
		Subsystem:  "histogram_subsystem",
		Name:       "histogram_name",
		Help:       "histogram_help",
		LabelNames: []string{"http", "https"},
	}

	histogram := p.NewHistogram(histogramOpts) // 默认的 buckets 是 []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}
	histogram.With("http", "www.baidu.com", "https", "github.com").Observe(0.8) // 大于 0.8 的 bucket 只有 {1, 2.5, 5, 10}
	histogram.With("http", "twitter", "https", "csdn.com").Observe(0.4) // 大于 0.4 的 bucket 有 {0.5, 1, 2.5, 5, 10}

	resp, err = client.Get(fmt.Sprintf("http://%s/metrics", server.Listener.Addr().String()))
	require.NoError(t, err)
	defer resp.Body.Close()

	bz, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	t.Log(string(bz))
}
