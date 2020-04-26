package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type options struct {
	port           int
	upstream       string
	upstreamParsed *url.URL

	logLevel string
}

func (o *options) validate() error {
	level, err := logrus.ParseLevel(o.logLevel)
	if err != nil {
		return fmt.Errorf("invalid log level specified: %v", err)
	}
	logrus.SetLevel(level)

	upstreamURL, err := url.Parse(o.upstream)
	if err != nil {
		return fmt.Errorf("failed to parse upstream URL: %v", err)
	}
	o.upstreamParsed = upstreamURL
	return nil
}

type defaultFieldsFormatter struct {
	WrappedFormatter logrus.Formatter
	DefaultFields    logrus.Fields
	PrintLineNumber  bool
}

// Format implements logrus.Formatter's Format.
func (f *defaultFieldsFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields, len(entry.Data)+len(f.DefaultFields))
	for k, v := range f.DefaultFields {
		data[k] = v
	}
	for k, v := range entry.Data {
		data[k] = v
	}
	return f.WrappedFormatter.Format(&logrus.Entry{
		Logger:  entry.Logger,
		Data:    data,
		Time:    entry.Time,
		Level:   entry.Level,
		Message: entry.Message,
		Caller:  entry.Caller,
	})
}

// type proxyTransport struct {
// 	http.RoundTripper
// }

// func (t *proxyTransport) RoundTrip(request *http.Request) (*http.Response, error) {
// 	response, err := t.RoundTripper.RoundTrip(request)
// 	body, err := httputil.DumpResponse(response, true)
// 	if err != nil {
// 		return nil, err
// 	}

// 	logrus.Info(string(body))

// 	return response, err
// }

func flagOptions() *options {
	o := &options{}
	flag.IntVar(&o.port, "port", 8080, "Port to listen on.")
	flag.StringVar(&o.upstream, "upstream", "https://hooks.slack.com", "Scheme, host, and base path of reverse proxy upstream.")
	flag.StringVar(&o.logLevel, "log-level", "debug", fmt.Sprintf("Log level is one of %v.", logrus.AllLevels))
	return o
}

func initLogrus() {
	formatter := defaultFieldsFormatter{
		PrintLineNumber:  true,
		DefaultFields:    logrus.Fields{"component": "rproxy"},
		WrappedFormatter: &logrus.JSONFormatter{},
	}

	logrus.SetFormatter(&formatter)
	logrus.SetReportCaller(formatter.PrintLineNumber)
	logrus.SetOutput(os.Stdout)
}

func main() {
	initLogrus()

	o := flagOptions()
	flag.Parse()
	if err := o.validate(); err != nil {
		logrus.WithError(err).Fatal("Invalid arguments.")
	}

	proxy := newReverseProxy(o.upstreamParsed, 30*time.Second)
	server := &http.Server{Addr: ":" + strconv.Itoa(o.port), Handler: proxy}

	logrus.Info("Server started.")

	logrus.WithError(server.ListenAndServe()).Info("Server exited.")
}

func newReverseProxy(upstreamURL *url.URL, timeout time.Duration) http.Handler {
	proxy := httputil.NewSingleHostReverseProxy(upstreamURL)
	// Wrap the director to change the upstream request 'Host' header to the
	// target host.
	director := proxy.Director
	proxy.Director = func(req *http.Request) {
		director(req)
		req.Host = req.URL.Host
	}
	// proxy.Transport = &proxyTransport{http.DefaultTransport}
	return http.TimeoutHandler(proxy, timeout, fmt.Sprintf("rproxy timed out after %v", timeout))
}
