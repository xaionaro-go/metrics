package metrics

import (
	"sync/atomic"
)

type metricCommonFloat64 struct {
	metricCommon
	modifyCounter uint64
	valuePtr      AtomicFloat64Interface
}

func (m *metricCommonFloat64) init(parent Metric, key string, tags AnyTags) {
	value := AtomicFloat64(0)
	m.valuePtr = &value
	m.metricCommon.init(parent, key, tags, func() bool { return m.wasUseless() })
}

func (m *metricCommonFloat64) Add(delta float64) float64 {
	if m == nil {
		return 0
	}
	if m.valuePtr == nil {
		return 0
	}
	r := m.valuePtr.Add(delta)
	atomic.AddUint64(&m.modifyCounter, 1)
	return r
}

func (m *metricCommonFloat64) Set(newValue float64) {
	if m == nil {
		return
	}
	if m.valuePtr == nil {
		return
	}
	m.valuePtr.Set(newValue)
	atomic.AddUint64(&m.modifyCounter, 1)
}

func (m *metricCommonFloat64) Get() float64 {
	if m == nil {
		return 0
	}
	if m.valuePtr == nil {
		return 0
	}
	return m.valuePtr.Get()
}

func (m *metricCommonFloat64) GetFloat64() float64 {
	return m.Get()
}

func (m *metricCommonFloat64) getModifyCounterDiffFlush() uint64 {
	if m == nil {
		return 0
	}
	return atomic.SwapUint64(&m.modifyCounter, 0)
}

func (w *metricCommonFloat64) SetValuePointer(newValuePtr *float64) {
	if w == nil {
		return
	}
	w.valuePtr = &AtomicFloat64Ptr{Pointer: newValuePtr}
}

func (m *metricCommonFloat64) Send(sender Sender) {
	if sender == nil {
		return
	}
	sender.SendFloat64(m.parent, string(m.storageKey), m.Get())
}

func (w *metricCommonFloat64) wasUseless() bool {
	return w.getModifyCounterDiffFlush() == 0
}
