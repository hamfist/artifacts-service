package server

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
)

// LoggerMiddleware is a middleware handler that logs the request as it goes in and the response as it goes out.
type LoggerMiddleware struct {
	// Logger is the log.Logger instance used to log messages with the Logger middleware
	Logger *logrus.Logger
}

// NewLoggerMiddleware returns a new Logger instance
func NewLoggerMiddleware() *LoggerMiddleware {
	return &LoggerMiddleware{
		Logger: logrus.New(),
	}
}

func (l *LoggerMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	l.Logger.WithFields(logrus.Fields{
		"method":  r.Method,
		"request": r.RequestURI,
		"remote":  r.RemoteAddr,
	}).Info("started handling request")

	next(rw, r)

	res := rw.(negroni.ResponseWriter)
	l.Logger.WithFields(logrus.Fields{
		"status":      res.Status(),
		"method":      r.Method,
		"request":     r.RequestURI,
		"remote":      r.RemoteAddr,
		"text_status": http.StatusText(res.Status()),
		"time":        time.Since(start),
	}).Info("completed handling request")
}
