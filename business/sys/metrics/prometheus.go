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

type LabelValuesFunc func() []string

type Histogram struct {
	hist       *prometheus.HistogramVec
	valuesFunc LabelValuesFunc
	begin      time.Time
}

func (h *Histogram) Send() {
	took := time.Since(h.begin)
	lvs := h.valuesFunc()
	h.hist.WithLabelValues(lvs...).Observe(took.Seconds())
}

// HttpHistogram creates a new Histogram responsible for observing http requests
func HttpHistogram(r *http.Request, v *web.CtxValues) *Histogram {
	return &Histogram{
		hist:       httpHist,
		valuesFunc: HTTPLabelValues(r, v),
		begin:      time.Now(),
	}
}

func HTTPLabelValues(r *http.Request, v *web.CtxValues) LabelValuesFunc {
	return func() []string {
		handler := r.URL.Path
		method := r.Method
		code := strconv.Itoa(v.StatusCode)

		return []string{handler, method, code}
	}
}

// DbHistogram creates a new Histogram responsible for observing db calls
func DbHistogram(table, operation string) *Histogram {
	return &Histogram{
		hist:       dbHist,
		valuesFunc: DBLabelValues(table, operation),
		begin:      time.Now(),
	}
}

func DBLabelValues(table, operation string) LabelValuesFunc {
	return func() []string {
		return []string{table, operation}
	}
}
