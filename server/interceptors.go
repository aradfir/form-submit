package server

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"math/rand"
	"time"
)

func loggerServerInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	requestId := rand.Uint64()
	log.WithFields(log.Fields{"request ID": requestId}).Info("Request  started")
	start := time.Now()
	h, err := handler(ctx, req)
	log.WithFields(log.Fields{
		"request ID": requestId,
		"method":     info.FullMethod,
		"duration":   time.Since(start),
		"error":      err,
	}).Info("Request finished")
	return h, err

}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "FormSubmit_processed_ops_total",
		Help: "The total number of processed events",
	})
	opsSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "FormSubmit_processed_ops_success",
		Help: "Number of successful events",
	})
	opsError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "FormSubmit_processed_ops_error",
		Help: "Number of requests that have errors",
	})
	responseTime = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "FormSubmit_response_time",
		Help: "The response time of requests",
	})
)

func instrumentalizationInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	opsProcessed.Inc()
	start := time.Now()
	h, err := handler(ctx, req)
	end := time.Since(start).Seconds()
	responseTime.Observe(end)
	if err != nil {
		opsError.Inc()
	} else {
		opsSuccess.Inc()
	}
	return h, err
}
