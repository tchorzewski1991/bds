package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tchorzewski1991/bds/base/web"
	"net/http"
	"strconv"
	"time"
)

var dbHist = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "db_call_duration",
	Help:    "The duration of db calls.",
	Buckets: prometheus.DefBuckets,
}, []string{"table", "operation"})

var httpHist = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "http_request_duration",
	Help:    "The duration of http requests.",
	Buckets: prometheus.DefBuckets,
}, []string{"handler", "method", "code"})

func init() {
	prometheus.Register(dbHist)
	prometheus.Register(httpHist)
}

type Histogram struct {
	hist  *prometheus.HistogramVec
	lvr   labelValuesReader
	begin time.Time
}

func (h *Histogram) Send() {
	took := time.Since(h.begin)
	lvs := h.lvr.Read()
	h.hist.WithLabelValues(lvs...).Observe(took.Seconds())
}

// HttpHistogram creates a new Histogram responsible for observing http requests
func HttpHistogram(r *http.Request, v *web.CtxValues) *Histogram {
	return &Histogram{
		hist: httpHist,
		lvr: &httpLabelValues{
			r: r,
			v: v,
		},
		begin: time.Now(),
	}
}

// DbHistogram creates a new Histogram responsible for observing db calls
func DbHistogram(table, operation string) *Histogram {
	return &Histogram{
		hist: dbHist,
		lvr: &dbLabelValues{
			table:     table,
			operation: operation,
		},
		begin: time.Now(),
	}
}

// private

type labelValuesReader interface {
	Read() []string
}

type httpLabelValues struct {
	r *http.Request
	v *web.CtxValues
}

func (v *httpLabelValues) Read() []string {
	handler := v.r.URL.Path
	method := v.r.Method
	code := strconv.Itoa(v.v.StatusCode)

	return []string{handler, method, code}
}

type dbLabelValues struct {
	table     string
	operation string
}

func (v *dbLabelValues) Read() []string {
	return []string{v.table, v.operation}
}
