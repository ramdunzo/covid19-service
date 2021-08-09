package newrelic

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	newrelic "github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrlogrus"
	"github.com/spf13/viper"
)

type Metric struct {
	Txn newrelic.Transaction
}

type MetricClosure struct {
	End func()
}

//Usages defer metrix.StartExternalSegment().End()
func (metric *Metric) StartExternalSegment(request *http.Request) *MetricClosure {

	s2 := newrelic.StartExternalSegment(metric.Txn, request)
	return &MetricClosure{
		End: func() {
			s2.End()
		},
	}
}

//Usages defer metrix.StartSegment(segmentName).End()
func (metric *Metric) StartSegment(name string) *MetricClosure {

	s := newrelic.StartSegment(metric.Txn, name)
	return &MetricClosure{
		End: func() {
			s.End()
		},
	}
}

// usage defer metrics.StartEsSegment(indexName, operation, param).End()
func (Metric *Metric) StartEsSegment(index, operation, param string) *MetricClosure {
	s := newrelic.DatastoreSegment{
		Product:            newrelic.DatastoreElasticsearch,
		Collection:         index,
		Operation:          operation,
		ParameterizedQuery: param,
	}
	s.StartTime = newrelic.StartSegmentNow(Metric.Txn)
	return &MetricClosure{
		End: func() {
			s.End()
		},
	}
}

//Usages defer metrix.StartSegment(segmentName).End()
func (metric *Metric) StartPostgreSegment(table, operation string) *MetricClosure {

	s := newrelic.DatastoreSegment{

		Product:    newrelic.DatastorePostgres,
		Collection: table,
		Operation:  operation,
	}
	s.StartTime = newrelic.StartSegmentNow(metric.Txn)
	return &MetricClosure{
		End: func() {
			s.End()
		},
	}
}

func (metric *Metric) StartRedisSegment(table, operation string) *MetricClosure {

	s := newrelic.DatastoreSegment{

		Product:    newrelic.DatastoreRedis,
		Collection: table,
		Operation:  operation,
	}
	s.StartTime = newrelic.StartSegmentNow(metric.Txn)
	return &MetricClosure{
		End: func() {
			s.End()
		},
	}
}

func (metric *Metric) StartMongoSegment(table, operation string) *MetricClosure {

	s := newrelic.DatastoreSegment{

		Product:    newrelic.DatastoreMongoDB,
		Collection: table,
		Operation:  operation,
	}
	s.StartTime = newrelic.StartSegmentNow(metric.Txn)
	return &MetricClosure{
		End: func() {
			s.End()
		},
	}
}

var NEW_RELIC_APP_NAME newrelic.Application

func SetUpNewRelic() {

	viper.BindEnv("NR_APPNAME")
	NEW_RELIC_APP_NAME = setUpNewRelic(viper.GetString("NR_APPNAME"))
}

func SetUpNewRelicWithCustomAppName(appName string) newrelic.Application {

	return setUpNewRelic(appName)
}

func setUpNewRelic(appName string) newrelic.Application {

	viper.BindEnv("ENV")
	viper.BindEnv("NEWRELICLICENSE")

	config := newrelic.NewConfig(appName, viper.GetString("NEWRELICLICENSE"))
	ignoreStatusCodesList := getIgnoredStatusCodesList()
	config.ErrorCollector.IgnoreStatusCodes = append(config.ErrorCollector.IgnoreStatusCodes, ignoreStatusCodesList...)
	if strings.ToLower(viper.GetString("newrelic.log_enabled")) == "true" {
		config.Logger = nrlogrus.StandardLogger()
	}
	if strings.ToLower(viper.GetString("newrelic.enabled")) != "true" {
		config.Enabled = false
	} else {
		config.Enabled = true
		config.DistributedTracer.Enabled = true
	}
	app, err := newrelic.NewApplication(config)
	if err != nil {
		panic(err)
	}
	return app
}

func getIgnoredStatusCodesList() []int {
	list := viper.GetStringSlice("ignored_status_codes")
	ignoredStatusCodes := []int{}
	for _, v := range list {
		i, err := strconv.Atoi(v)
		if err != nil {
			panic("invalid status code.")
		}
		ignoredStatusCodes = append(ignoredStatusCodes, i)
	}
	return ignoredStatusCodes
}

func Instrument(v0mux *mux.Router, apiPath string, usersHandler func(w http.ResponseWriter, req *http.Request)) *mux.Route {

	return v0mux.HandleFunc(newrelic.WrapHandleFunc(NEW_RELIC_APP_NAME, apiPath, usersHandler))
}

func InstrumentWithApiPath(v0mux *mux.Router, apiPath string, usersHandler func(w http.ResponseWriter, req *http.Request)) *mux.Route {

	return v0mux.HandleFunc(wrapHandleFuncCustom(NEW_RELIC_APP_NAME, apiPath, usersHandler))
}

func InstrumentWithMatcher(f func(request *http.Request, match *mux.RouteMatch) bool, v0mux *mux.Router,
	customNewRelicApp newrelic.Application, handler http.Handler) *mux.Route {

	_, h := wrapHandlerCustom(customNewRelicApp, "", handler)
	return v0mux.MatcherFunc(f).Handler(h)
}

func InstrumentWithHandler(v0mux *mux.Router, apiPath string, customNewRelicApp newrelic.Application,
	handler http.Handler) *mux.Route {

	_, h := wrapHandlerCustom(customNewRelicApp, "", handler)
	return v0mux.Handle(apiPath, h)
}

func StartTxn(txnName string) *Metric {

	return &Metric{NEW_RELIC_APP_NAME.StartTransaction(txnName, nil, nil)}
}

func EndTxn(metric *Metric) {
	metric.Txn.End()
}

// WrapHandleFunc serves the same purpose as WrapHandle for functions registered
// with ServeMux.HandleFunc.
func wrapHandleFuncCustom(app newrelic.Application, pattern string, handler func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request)) {
	p, h := wrapHandleCustom(app, pattern, http.HandlerFunc(handler))
	return p, func(w http.ResponseWriter, r *http.Request) { h.ServeHTTP(w, r) }
}

func wrapHandlerCustom(app newrelic.Application, pattern string, handler http.Handler) (string, http.Handler) {
	p, h := wrapHandleCustom(app, pattern, handler)
	return p, h
}

func wrapHandleCustom(app newrelic.Application, pattern string, handler http.Handler) (string, http.Handler) {
	if app == nil {
		return pattern, handler
	}
	return pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		txn := app.StartTransaction(r.URL.Path, w, r)
		defer txn.End()

		r = newrelic.RequestWithTransactionContext(r, txn)

		handler.ServeHTTP(txn, r)
	})
}

func (m *Metric) AddAttribute(key string, val interface{}) {
	if m != nil {
		if m.Txn != nil {
			m.Txn.AddAttribute(key, val)
		}
	}
}
