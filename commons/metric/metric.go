// Copyright 2014 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metric

import (
	"encoding/json"
	"errors"
	"math"
	"sync"
	"sync/atomic"
)

type counter struct {
	// valBits contains the bits of the represented float64 value, while
	// valInt stores values that are exact integers. Both have to go first
	// in the struct to guarantee alignment for atomic operations.
	// http://golang.org/pkg/sync/atomic/#pkg-note-BUG
	valBits uint64
	valInt  uint64
}

func NewCounter() *counter {
	return &counter{
		valInt:  0,
		valBits: 0,
	}
}

func strjson(x interface{}) string {
	b, _ := json.Marshal(x)
	return string(b)
}
func (c *counter) Val() float64 {
	fval := math.Float64frombits(atomic.LoadUint64(&c.valBits))
	ival := atomic.LoadUint64(&c.valInt)
	val := fval + float64(ival)
	return val
}

func (c *counter) String() string { return strjson(c) }
func (c *counter) Reset() {
	atomic.StoreUint64(&c.valBits, 0)
	atomic.StoreUint64(&c.valInt, 0)
}
func (c *counter) Add(v float64) {
	if v < 0 {
		panic(errors.New("counter cannot decrease in value"))
	}
	ival := uint64(v)
	if float64(ival) == v {
		atomic.AddUint64(&c.valInt, ival)
		return
	}
	for {
		oldBits := atomic.LoadUint64(&c.valBits)
		newBits := math.Float64bits(math.Float64frombits(oldBits) + v)
		if atomic.CompareAndSwapUint64(&c.valBits, oldBits, newBits) {
			return
		}
	}
}
func (c *counter) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type  string  `json:"type"`
		Count float64 `json:"count"`
	}{"c", c.Val()})
}

type bin struct {
	value float64
	count float64
}

type histogram struct {
	sync.Mutex
	bins    []bin
	total   uint64
	maxBins int
}

func NewHistogram() *histogram {
	return &histogram{
		bins:    []bin{},
		maxBins: 1000,
	}
}

func (h *histogram) Add(n float64) {
	h.Lock()
	defer h.Unlock()
	defer h.trim()
	h.total++
	// sort 动态的从小到大排序
	newbin := bin{value: n, count: 1}
	for i := range h.bins {
		if h.bins[i].value > n {
			h.bins = append(h.bins[:i], append([]bin{newbin}, h.bins[i:]...)...)
			return
		}
	}
	h.bins = append(h.bins, newbin)
}

func (h *histogram) MarshalJSON() ([]byte, error) {
	h.Lock()
	defer h.Unlock()
	return json.Marshal(struct {
		Type string  `json:"type"`
		P50  float64 `json:"p50"`
		P90  float64 `json:"p90"`
		P99  float64 `json:"p99"`
		P999 float64 `json:"p99.9"`
	}{"h", h.quantile(0.5), h.quantile(0.9), h.quantile(0.99), h.quantile(0.999)})
}

func (h *histogram) trim() {
	for len(h.bins) > h.maxBins {
		d := float64(0)
		i := 0
		for j := 1; j < len(h.bins); j++ {
			if dv := h.bins[j].value - h.bins[j-1].value; dv < d || j == 1 {
				d = dv
				i = j
			}
		}
		count := h.bins[i-1].count + h.bins[i].count
		merged := bin{
			value: (h.bins[i-1].value*h.bins[i-1].count + h.bins[i].value*h.bins[i].count) / count,
			count: count,
		}
		h.bins = append(h.bins[:i-1], h.bins[i:]...)
		h.bins[i-1] = merged
	}
}

func (h *histogram) quantile(q float64) float64 {
	count := q * float64(h.total)
	for i := range h.bins {
		count -= h.bins[i].count
		if count <= 0 {
			return h.bins[i].value
		}
	}
	return 0
}

type gauge struct {
	sync.Mutex
	sum   float64
	min   float64
	max   float64
	count int
}

func NewGauge() *gauge {
	return &gauge{}
}

func (g *gauge) String() string { return strjson(g) }
func (g *gauge) Reset() {
	g.Lock()
	defer g.Unlock()
	g.count, g.sum, g.min, g.max = 0, 0, 0, 0
}

func (g *gauge) Add(n float64) {
	g.Lock()
	defer g.Unlock()
	if n < g.min || g.count == 0 {
		g.min = n
	}
	if n > g.max || g.count == 0 {
		g.max = n
	}
	g.sum += n
	g.count++
}

func (g *gauge) MarshalJSON() ([]byte, error) {
	g.Lock()
	defer g.Unlock()
	return json.Marshal(struct {
		Type string  `json:"type"`
		Mean float64 `json:"mean"`
		Min  float64 `json:"min"`
		Max  float64 `json:"max"`
	}{"g", g.mean(), g.min, g.max})
}
func (g *gauge) mean() float64 {
	if g.count == 0 {
		return 0
	}
	return g.sum / float64(g.count)
}
