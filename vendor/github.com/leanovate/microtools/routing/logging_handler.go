package routing

import (
	"net/http"
	"time"

	"github.com/leanovate/microtools/logging"
)

var flowIDHeaderNames = []string{"X-Flow-Id"}

type loggingResponseWriter struct {
	status        int
	responseBytes int
	underlying    http.ResponseWriter
}

func (l *loggingResponseWriter) Header() http.Header {
	return l.underlying.Header()
}

func (l *loggingResponseWriter) Write(bytes []byte) (int, error) {
	n, err := l.underlying.Write(bytes)
	if err != nil {
		l.responseBytes += n
	}
	return n, err
}

func (l *loggingResponseWriter) WriteHeader(status int) {
	l.status = status
	l.underlying.WriteHeader(status)
}

// LoggingHandler is a http.Handler that logs and delegates all requests.
type LoggingHandler struct {
	delegate http.Handler
	logger   logging.Logger
}

// NewLoggingHandler create a new LoggingHandler.
func NewLoggingHandler(delegate http.Handler, logger logging.Logger) http.Handler {
	return &LoggingHandler{
		delegate: delegate,
		logger:   logger,
	}
}

func (l *LoggingHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	start := time.Now()
	loggingResp := &loggingResponseWriter{
		status:     -1,
		underlying: resp,
	}
	l.delegate.ServeHTTP(loggingResp, req)
	elapsed := time.Since(start)
	flowID := flowIDFromHeaders(req)
	log := l.logger.WithContext(map[string]interface{}{
		"method":  req.Method,
		"uri":     req.RequestURI,
		"status":  loggingResp.status,
		"time":    elapsed.String(),
		"millis":  float64(elapsed.Nanoseconds()) / 1000000.0,
		"bytes":   loggingResp.responseBytes,
		"flow_id": flowID,
	})
	if loggingResp.status < 300 {
		log.Info("Request: Success")
	} else if loggingResp.status < 400 {
		log.Info("Request: Redirect")
	} else if loggingResp.status < 500 {
		log.Warn("Request: Client error")
	} else {
		log.Error("Request: Server error")
	}
}

func flowIDFromHeaders(req *http.Request) string {
	for _, headerName := range flowIDHeaderNames {
		if flowID := req.Header.Get(headerName); flowID != "" {
			return flowID
		}
	}
	return ""
}
