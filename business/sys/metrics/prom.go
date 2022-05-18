package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tchorzewski1991/bds/base/web"
	"net/http"
	"strconv"
	"time"
)

var httpHist = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "http_request_duration",
	Help:    "The duration of http requests.",
	Buckets: prometheus.DefBuckets,
}, []string{"handler", "method", "code"})

var httpCount = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "http_request_total",
	Help: "The number of http requests.",
}, []string{"handler", "method", "code"})

var dbHist = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "db_call_duration",
	Help:    "The duration of db calls.",
	Buckets: prometheus.DefBuckets,
}, []string{"operation"})

var dbCount = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "db_call_total",
	Help: "The number of db calls.",
}, []string{"operation"})

func init() {
	prometheus.Register(httpHist)
	prometheus.Register(httpCount)
	prometheus.Register(dbHist)
	prometheus.Register(dbCount)
}

type HttpHistogram struct {
	hist  *prometheus.HistogramVec
	req   *http.Request
	val   *web.CtxValues
	begin time.Time
}

func NewHttpHistogram(req *http.Request, val *web.CtxValues) *HttpHistogram {
	return &HttpHistogram{
		begin: time.Now(),
		hist:  httpHist,
		req:   req,
		val:   val,
	}
}

func (o *HttpHistogram) Send() {
	handler := o.req.URL.Path
	method := o.req.Method
	code := strconv.Itoa(o.val.StatusCode)
	took := time.Since(o.begin)

	o.hist.WithLabelValues(handler, method, code).Observe(took.Seconds())
}

type DbHistogram struct {
	begin     time.Time
	hist      *prometheus.HistogramVec
	operation string
}

func NewDBHistogram(operation string) *DbHistogram {
	return &DbHistogram{
		begin:     time.Now(),
		hist:      dbHist,
		operation: operation,
	}
}

func (m *DbHistogram) Send() {
	m.hist.WithLabelValues(m.operation).Observe(time.Since(m.begin).Seconds())
}
