package metrics

import (
	"sync/atomic"
	"unsafe"
)

type MetricGaugeInt64 struct {
	metricCommonInt64
}

func newMetricGaugeInt64(key string, tags AnyTags) *MetricGaugeInt64 {
	metric := metricGaugeInt64Pool.Get().(*MetricGaugeInt64)
	metric.init(key, tags)
	return metric
}

func (m *MetricGaugeInt64) init(key string, tags AnyTags) {
	m.metricCommonInt64.init(m, key, tags)
}

func GaugeInt64(key string, tags AnyTags) *MetricGaugeInt64 {
	if IsDisabled() {
		return (*MetricGaugeInt64)(nil)
	}

	m := registry.Get(TypeGaugeInt64, key, tags)
	if m != nil {
		return m.(*MetricGaugeInt64)
	}

	return newMetricGaugeInt64(key, tags)
}

func (m *MetricGaugeInt64) GetType() Type {
	return TypeGaugeInt64
}

func (m *MetricGaugeInt64) Decrement() int64 {
	if m == nil {
		return 0
	}
	if m.valuePtr == nil {
		atomic.StorePointer((*unsafe.Pointer)((unsafe.Pointer)(&m.valuePtr)), (unsafe.Pointer)(&[]int64{0}[0]))
	}
	return atomic.AddInt64(m.valuePtr, -1)
}
