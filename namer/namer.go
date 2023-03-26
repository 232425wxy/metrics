package namer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/232425wxy/metrics"
)

type Namer struct {
	namespace  string
	subsystem  string
	name       string
	nameFormat string
	labelNames map[string]struct{}
}

// formatRegexp
// 表达式[:alnum:]表示匹配任何英文大小写字母和数字
// 表达式([#?[:alnum:]_]+)表示匹配至少一个井号("#")或者问号("?")或者字母，再或者是下划线("_")，方括号内的特殊字符会被当成普通字符来匹配，例如此处的问号("?")。
// 表达式%{([#?[:alnum:]_]+)}匹配以"%{"开头，然后以"}"结尾的字符串，中间的内容必须满足表达式([#?[:alnum:]_]+)的匹配规则。
var formatRegexp = regexp.MustCompile(`%{([#?[:alnum:]_]+)}`)

// invalidLabelValueRegexp 匹配字符串里的空格(" ")、点号(".")、冒号(":")以及竖杠("|")，方括号内的特殊字符会被当成普通字符来匹配，例如此处的竖杠("|")。
var invalidLabelValueRegexp = regexp.MustCompile(`[.|:\s]`)

func NewCounterNamer(opts metrics.CounterOpts) *Namer {
	return &Namer{
		namespace:  opts.Namespace,
		subsystem:  opts.Subsystem,
		name:       opts.Name,
		nameFormat: opts.StatsdFormat,
		labelNames: sliceToSet(opts.LabelNames),
	}
}

func NewGaugeNamer(opts metrics.GaugeOpts) *Namer {
	return &Namer{
		namespace:  opts.Namespace,
		subsystem:  opts.Subsystem,
		name:       opts.Name,
		nameFormat: opts.StatsdFormat,
		labelNames: sliceToSet(opts.LabelNames),
	}
}

func NewHistogramNamer(opts metrics.HistogramOpts) *Namer {
	return &Namer{
		namespace:  opts.Namespace,
		subsystem:  opts.Subsystem,
		name:       opts.Name,
		nameFormat: opts.StatsdFormat,
		labelNames: sliceToSet(opts.LabelNames),
	}
}

// FullyQualifiedName 直译过来就是：`完全合格的名字`。
// return namespace.subsystem.name or namespace.name or subsystem.name or name
func (n *Namer) FullyQualifiedName() string {
	switch {
	case n.namespace != "" && n.subsystem != "":
		return strings.Join([]string{n.namespace, n.subsystem, n.name}, ".")
	case n.namespace != "":
		return strings.Join([]string{n.namespace, n.name}, ".")
	case n.subsystem != "":
		return strings.Join([]string{n.subsystem, n.name}, ".")
	default:
		return n.name
	}
}

// Format 传入的labelValues里的label必须是在Namer.labelNames里存在的！
func (n *Namer) Format(labelValues ...string) string {
	labels := n.labelsToMap(labelValues)

	cursor := 0
	var segments []string

	matches := formatRegexp.FindAllStringSubmatchIndex(n.nameFormat, -1)
	for _, m := range matches {
		start, end := m[0], m[1]           // 匹配 %{xxx}。
		labelStart, labelEnd := m[2], m[3] // 匹配 %{xxx} 里的 xxx。

		if start > cursor {
			segments = append(segments, n.nameFormat[cursor:start]) // 这里将 yyy%{xxx} 字符串里的 yyy 添加到segments里。
		}

		key := n.nameFormat[labelStart:labelEnd] // 获取 %{xxx} 里的 xxx。
		var value string
		switch key {
		case "#namespace":
			value = n.namespace
		case "#subsystem":
			value = n.subsystem
		case "#name":
			value = n.name
		case "#fqname":
			value = n.FullyQualifiedName()
		default:
			var ok bool
			value, ok = labels[key]
			if !ok {
				panic(fmt.Sprintf("invalid label in name format: %s", key))
			}
			value = invalidLabelValueRegexp.ReplaceAllString(value, "_") // 将 value 里的所有空格(" ")、点号(".")、冒号(":")以及竖杠("|")替换为下划线("_")。
		}
		segments = append(segments, value)
		cursor = end
	}

	if cursor != len(n.nameFormat) {
		segments = append(segments, n.nameFormat[cursor:])
	}

	return strings.Join(segments, "")
}

// labelsToMap 要求给定的labels必须在Namer的labelNames里存在。
func (n *Namer) labelsToMap(labelValues []string) map[string]string {
	labels := map[string]string{}
	for i := 0; i < len(labelValues); i += 2 {
		key := labelValues[i]
		n.validateKey(key)
		if i == len(labelValues)-1 {
			labels[key] = "unknown"
		} else {
			labels[key] = labelValues[i+1]
		}
	}
	return labels
}

func (n *Namer) validateKey(label string) {
	if _, ok := n.labelNames[label]; !ok {
		panic("invalid label name: " + label)
	}
}

// sliceToSet 将切片转换为map，空结构体不占用任何存储空间。
func sliceToSet(slice []string) map[string]struct{} {
	labelSet := make(map[string]struct{})
	for _, s := range slice {
		labelSet[s] = struct{}{}
	}
	return labelSet
}
