package prometheus

import (
	"github.com/232425wxy/metrics"
	kitmetrics "github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	prom "github.com/prometheus/client_golang/prometheus"
)

type Provider struct{}

type Counter struct {
	kitmetrics.Counter
}

func (c *Counter) With(labelValues ...string) metrics.Counter {
	return &Counter{Counter: c.Counter.With(labelValues...)}
}

type Gauge struct {
	kitmetrics.Gauge
}

func (g *Gauge) With(labelValues ...string) metrics.Gauge {
	return &Gauge{Gauge: g.Gauge.With(labelValues...)}
}

type Histogram struct {
	kitmetrics.Histogram
}

func (h *Histogram) With(labelValues ...string) metrics.Histogram {
	return &Histogram{Histogram: h.Histogram.With(labelValues...)}
}

func (p *Provider) NewCounter(opts metrics.CounterOpts) metrics.Counter {
	return &Counter{
		Counter: prometheus.NewCounterFrom(
			prom.CounterOpts{
				Namespace: opts.Namespace,
				Subsystem: opts.Subsystem,
				Name:      opts.Name,
				Help:      opts.Help,
			},
			opts.LabelNames,
		),
	}
}

func (p *Provider) NewGauge(opts metrics.GaugeOpts) metrics.Gauge {
	return &Gauge{
		Gauge: prometheus.NewGaugeFrom(
			prom.GaugeOpts{
				Namespace: opts.Namespace,
				Subsystem: opts.Subsystem,
				Name:      opts.Name,
				Help:      opts.Help,
			},
			opts.LabelNames,
		),
	}
}

func (p *Provider) NewHistogram(opts metrics.HistogramOpts) metrics.Histogram {
	return &Histogram{
		Histogram: prometheus.NewHistogramFrom(
			prom.HistogramOpts{
				Namespace: opts.Namespace,
				Subsystem: opts.Subsystem,
				Name:      opts.Name,
				Help:      opts.Help,
				Buckets:   opts.Buckets,
			},
			opts.LabelNames,
		),
	}
}
