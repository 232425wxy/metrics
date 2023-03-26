package statsd

import (
	"github.com/232425wxy/metrics"
	"github.com/232425wxy/metrics/namer"
	"github.com/go-kit/kit/metrics/statsd"
)

const defaultFormat = "%{#fqname}"

type Provider struct {
	Statsd *statsd.Statsd
}

type Counter struct {
	Counter        *statsd.Counter
	namer          *namer.Namer
	statsdProvider *statsd.Statsd
}

func (c *Counter) Add(delta float64) {
	if c.Counter == nil {
		panic("label values must be provided by calling With")
	}
	c.Counter.Add(delta)
}

func (c *Counter) With(labelValues ...string) metrics.Counter {
	name := c.namer.Format(labelValues...)
	return &Counter{Counter: c.statsdProvider.NewCounter(name, 1)}
}

type Gauge struct {
	Gauge          *statsd.Gauge
	namer          *namer.Namer
	statsdProvider *statsd.Statsd
}

func (g *Gauge) Add(delta float64) {
	if g.Gauge == nil {
		panic("label values must be provided by calling With")
	}
	g.Gauge.Add(delta)
}

func (g *Gauge) Set(value float64) {
	if g.Gauge == nil {
		panic("label values must be provided by calling With")
	}
	g.Gauge.Set(value)
}

func (g *Gauge) With(labelValues ...string) metrics.Gauge {
	name := g.namer.Format(labelValues...)
	return &Gauge{Gauge: g.statsdProvider.NewGauge(name)}
}

type Histogram struct {
	Timing *statsd.Timing
	namer *namer.Namer
	statsdProvider *statsd.Statsd
}

func (h *Histogram) With(labelValues ...string) metrics.Histogram {
	name := h.namer.Format(labelValues...)
	return &Histogram{Timing: h.statsdProvider.NewTiming(name, 1)}
}

func (h *Histogram) Observe(value float64) {
	if h.Timing == nil {
		panic("label values must be provided by calling With")
	}
	h.Timing.Observe(value)
}

// NewCounter 新建一个Counter，如果opts里的LabelNames是空的，则暂时不会创建statsd的Counter，
// 所以将来调用Add方法会因为空指针调用而panic，所以想避免该错误，则必须在调用Add方法之前，先调
// 用With方法。
// 注意：往With方法里传入的labelValues，所有的labels必须是LabelNames里所包含的，不然会panic。
func (p *Provider) NewCounter(opts metrics.CounterOpts) metrics.Counter {
	if opts.StatsdFormat == "" {
		opts.StatsdFormat = defaultFormat
	}
	counter := &Counter{
		namer:          namer.NewCounterNamer(opts),
		statsdProvider: p.Statsd,
	}

	if len(opts.LabelNames) == 0 {
		counter.Counter = p.Statsd.NewCounter(counter.namer.Format(), 1)
	}

	return counter
}

// NewGauge 新建一个Gauge，如果opts里的LabelNames是空的，则暂时不会创建statsd的Gauge，
// 所以将来调用Set或Add方法会因为空指针调用而panic，所以想避免该错误，则必须在调用Set或
// Add方法之前，先调用With方法。
// 注意：往With方法里传入的labelValues，所有的labels必须是LabelNames里所包含的，不然会panic。
func (p *Provider) NewGauge(opts metrics.GaugeOpts) metrics.Gauge {
	if opts.StatsdFormat == "" {
		opts.StatsdFormat = defaultFormat
	}
	gauge := &Gauge{
		namer:          namer.NewGaugeNamer(opts),
		statsdProvider: p.Statsd,
	}

	if len(opts.LabelNames) == 0 {
		gauge.Gauge = p.Statsd.NewGauge(gauge.namer.Format())
	}

	return gauge
}

// NewHistogram 新建一个Histogram，如果opts里的LabelNames是空的，则暂时不会创建statsd的Histogram，
// 所以将来调用Observe方法会因为空指针调用而panic，所以想避免该错误，则必须在调用Observe方法之前，先
// 调用With方法。
// 注意：往With方法里传入的labelValues，所有的labels必须是LabelNames里所包含的，不然会panic。
func (p *Provider) NewHistogram(opts metrics.HistogramOpts) metrics.Histogram {
	if opts.StatsdFormat == "" {
		opts.StatsdFormat = defaultFormat
	}
	histogram := &Histogram{
		namer:          namer.NewHistogramNamer(opts),
		statsdProvider: p.Statsd,
	}

	if len(opts.LabelNames) == 0 {
		histogram.Timing = p.Statsd.NewTiming(histogram.namer.Format(), 1)
	}

	return histogram
}