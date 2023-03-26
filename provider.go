package metrics

type Provider interface {
	NewCounter(CounterOpts) Counter
	NewGauge(GaugeOpts) Gauge
	NewHistogram(HistogramOpts) Histogram
}

// Counter 是一个累计类型的数据指标，它代表单调递增的计数器，其值只能增加，不能减少。
// 例如使用Counter来表示响应的HTTP请求数，这个数一定是不断增长的。
type Counter interface {
	With(labelValues ...string) Counter
	// 计数器增加值
	Add(delta float64)
}

type CounterOpts struct {
	Namespace    string
	Subsystem    string
	Name         string
	Help         string
	LabelNames   []string
	LabelHelp    map[string]string
	StatsdFormat string
}

// Gauge 是可以任意上下波动数值的指标类型，也就是说Gauge的值可增可减。
// 例如电脑CPU的使用率，可大可小。
type Gauge interface {
	With(labelValues ...string) Gauge
	Add(delta float64)
	Set(value float64)
}

type GaugeOpts struct {
	Namespace    string
	Subsystem    string
	Name         string
	Help         string
	LabelNames   []string
	LabelHelp    map[string]string
	StatsdFormat string
}

// Histogram 在Prometheus里是一种累积直方图，在弄懂什么是累积直方图前，先看一个例子：
//
// Example: 假设我们想监控某个应用在一段时间内的响应时间，最后监控到的样本的响应时间范围
// 为 0s~10s。现在我们将样本的值域划分为不同的区间，即不同的 bucket，每个 bucket 的宽度
// 是 0.2s。那么第一个 bucket 表示响应时间小于等于 0.2s 的请求数量，第二个 bucket 表示
// 响应时间大于 0.2s 小于等于 0.4s 的请求数量，以此类推。
//
// Prometheus 的 Histogram 与上面的区间划分方式是有差别的，它的划分方式如下：还假设每个
// bucket 的宽度是 0.2s，那么第一个 bucket 表示响应时间小于等于 0.2s 的请求数量，第二个
// bucket 表示响应时间小于等于 0.4s 的请求数量，以此类推。也就是说，每一个 bucket 的样本
// 包含了之前所有 bucket 的样本，所以叫累积直方图。
type Histogram interface {
	With(labelValues ...string) Histogram
	Observe(value float64)
}

type HistogramOpts struct {
	Namespace    string
	Subsystem    string
	Name         string
	Help         string
	Buckets      []float64
	LabelNames   []string
	LabelHelp    map[string]string
	StatsdFormat string
}
