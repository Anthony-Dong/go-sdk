package metric

import (
	"testing"
)

func Test_histogram_Add(t *testing.T) {
	h := NewHistogram()
	for x := 0; x < 10000; x++ {
		h.Add(float64(x))
	}
	json, _ := h.MarshalJSON()
	t.Log(string(json))
}
func TestCounter(t *testing.T) {
	t.Log(uint64(float64(1)))
	h := NewCounter()
	for x := 0; x < 10000; x++ {
		h.Add(1 + (1 / float64(x+1)))
	}
	t.Log(h.String())
}

func TestName(t *testing.T) {
	g := NewGauge()
	for x := 0; x < 10000; x++ {
		g.Add(float64(x))
	}
	t.Log(g.String())
	t.Log(uint64(g.sum))
	t.Log(uint64(g.count))
}
