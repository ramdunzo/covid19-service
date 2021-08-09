package common_dtos

import (
	"net/http"
	"strings"

	"github.com/dunzoit/dunzo_commons/go_commons/common-utils/common_logger"
	"github.com/dunzoit/dunzo_commons/go_commons/common-utils/id_generator"
	"github.com/dunzoit/dunzo_commons/go_commons/common-utils/service_commons/newrelic"
	newrelic2 "github.com/newrelic/go-agent"
	"github.com/sirupsen/logrus"
)

const CorrelationIdHeaderName = "CORRELATION-ID"

// Each of the fields embedded here is RO
// Unexpected behaviour might occur if the logger is updated inside a goroutine
type Context struct {
	Metric        *newrelic.Metric
	CorrelationId *string
	Data          *map[string]interface{}
	Logger        common_logger.Logger
}

func CreateContextFromRequest(request *http.Request) *Context {

	correlationId := GetCorrelationId(request)
	return &Context{Metric: &newrelic.Metric{Txn: newrelic2.FromContext(request.Context())}, CorrelationId: &correlationId}
}

func CreateLoggableContextFromRequest(request *http.Request) *Context {

	correlationId := GetCorrelationId(request)
	metric := &newrelic.Metric{
		Txn: newrelic2.FromContext(request.Context()),
	}
	metric.AddAttribute(CorrelationIdHeaderName, correlationId)
	return &Context{
		Metric:        metric,
		CorrelationId: &correlationId,
		Logger:        common_logger.CreateLoggerForRequest(request).With(common_logger.String(strings.ToLower(CorrelationIdHeaderName), correlationId)),
	}
}

func CreateLoggableContextWithNonWebTxn(txnName, correlationId string) (*Context, func()) {
	metric := newrelic.StartTxn(txnName)
	metric.AddAttribute(CorrelationIdHeaderName, correlationId)
	return &Context{
		Metric:        metric,
		CorrelationId: &correlationId,
		Logger:        common_logger.GlobalLogger().With(common_logger.String(strings.ToLower(CorrelationIdHeaderName), correlationId), common_logger.String("txnName", txnName)),
	}, func() {
		newrelic.EndTxn(metric)
	}
}

func GetCorrelationId(request *http.Request) string {

	correlationId := request.Header.Get(CorrelationIdHeaderName)
	if correlationId == "" {
		correlationId = id_generator.GetUniqId()
		logrus.Info("Creating new co-relation id : ", correlationId)
	}
	return correlationId
}

func SetCorrelationId(request *http.Request, ctx *Context) {

	if ctx == nil || ctx.CorrelationId == nil || *ctx.CorrelationId == "" {
		request.Header.Set(CorrelationIdHeaderName, id_generator.GetUniqId())
	} else {
		request.Header.Set(CorrelationIdHeaderName, *ctx.CorrelationId)
	}
	logrus.Info("Setting co-relation id as : ", *ctx.CorrelationId)
}
