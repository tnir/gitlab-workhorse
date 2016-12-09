package helper

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	responseLogger *log.Logger

	sessionsActive = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "gitlab_workhorse_http_sessions_active",
		Help: "Number of HTTP request-response cycles currently being handled by gitlab-workhorse.",
	})

	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gitlab_workhorse_http_requests_total",
			Help: "How many HTTP requests have been processed by gitlab-workhorse, partitioned by status code and HTTP method.",
		},
		[]string{"code", "method"},
	)

	cloneFetchRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gitlab_workhorse_git_clone_fetch_requests",
			Help: "How many Git clone/fetch requests for CI have been processed by gitlab-workhorse, partitioned by CI yes/no status.",
		},
		[]string{"ci"},
	)

	cloneFetchBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gitlab_workhorse_git_clone_fetch_bytes",
			Help: "How many Git clone/fetch bytes for CI have been send by gitlab-workhorse, partitioned by CI yes/no status.",
		},
		[]string{"ci"},
	)
)

func init() {
	SetCustomResponseLogger(os.Stderr)
	registerPrometheusMetrics()
}

func SetCustomResponseLogger(writer io.Writer) {
	responseLogger = log.New(writer, "", 0)
}

func registerPrometheusMetrics() {
	prometheus.MustRegister(sessionsActive)
	prometheus.MustRegister(requestsTotal)
	prometheus.MustRegister(cloneFetchRequests)
	prometheus.MustRegister(cloneFetchBytes)
}

type LoggingResponseWriter struct {
	rw      http.ResponseWriter
	status  int
	written int64
	started time.Time
}

func NewLoggingResponseWriter(rw http.ResponseWriter) LoggingResponseWriter {
	sessionsActive.Inc()
	return LoggingResponseWriter{
		rw:      rw,
		started: time.Now(),
	}
}

func (l *LoggingResponseWriter) Header() http.Header {
	return l.rw.Header()
}

func (l *LoggingResponseWriter) Write(data []byte) (n int, err error) {
	if l.status == 0 {
		l.WriteHeader(http.StatusOK)
	}
	n, err = l.rw.Write(data)
	l.written += int64(n)
	return
}

func (l *LoggingResponseWriter) WriteHeader(status int) {
	if l.status != 0 {
		return
	}

	l.status = status
	l.rw.WriteHeader(status)
}

func (l *LoggingResponseWriter) Log(r *http.Request) {
	duration := time.Since(l.started)
	responseLogger.Printf("%s %s - - [%s] %q %d %d %q %q %f\n",
		r.Host, r.RemoteAddr, l.started,
		fmt.Sprintf("%s %s %s", r.Method, r.RequestURI, r.Proto),
		l.status, l.written, r.Referer(), r.UserAgent(), duration.Seconds(),
	)

	l.countCloneFetchRequests(r)

	sessionsActive.Dec()
	requestsTotal.WithLabelValues(strconv.Itoa(l.status), r.Method).Inc()
}

func (l *LoggingResponseWriter) countCloneFetchRequests(r *http.Request) {
	if l.status == 401 || !strings.Contains(r.RequestURI, "/info/refs?service=git-upload-pack") {
		return
	}

	u, _, ok := r.BasicAuth()
	var forCi string
	if ok && u == "gitlab-ci-token" {
		forCi = "1"
	} else {
		forCi = "0"
	}

	cloneFetchRequests.WithLabelValues(forCi).Inc()
	cloneFetchBytes.WithLabelValues(forCi).Add(float64(l.written))
}
