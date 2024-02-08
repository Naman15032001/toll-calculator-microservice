package main

import (
	"time"

	"github.com/Naman15032001/tolling/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next Aggregator
}

type MetricsMiddleware struct {
	reqCounterAgg    prometheus.Counter
	reqCounterCalc   prometheus.Counter
	reqLatencyAgg    prometheus.Histogram
	reqLatencyCalc   prometheus.Histogram
	errorCounterAgg  prometheus.Counter
	errorCounterCalc prometheus.Counter
	next             Aggregator
}

func NewMetricsMiddleware(next Aggregator) *MetricsMiddleware {
	errorCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
		Name:      "aggregate",
	})

	errorCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
		Name:      "calculate",
	})
	reqCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "aggregate",
	})

	reqCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "calculate",
	})

	reqLatencyAgg := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "aggregate",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	reqLatencyCalc := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "calculate",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	return &MetricsMiddleware{
		reqCounterAgg:    reqCounterAgg,
		reqCounterCalc:   reqCounterCalc,
		reqLatencyAgg:    reqLatencyAgg,
		reqLatencyCalc:   reqLatencyCalc,
		errorCounterAgg:  errorCounterAgg,
		errorCounterCalc: errorCounterCalc,
		next:             next,
	}
}

func (m *MetricsMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		m.reqLatencyAgg.Observe(time.Since(start).Seconds())
		m.reqCounterAgg.Inc()
		if err != nil {
			m.errorCounterAgg.Inc()
		}
	}(time.Now())
	err = m.next.AggregateDistance(distance)
	return err
}

// func (m *MetricsMiddleware) pushMetrics(start time.Time , err error){
// 	m.reqCounter
// }

func (m *MetricsMiddleware) CalculateInvoice(obuid int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		m.reqLatencyCalc.Observe(time.Since(start).Seconds())
		m.reqCounterCalc.Inc()
		if err != nil {
			m.errorCounterCalc.Inc()
		}
	}(time.Now())
	inv, err = m.next.CalculateInvoice(obuid)
	return inv, err
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
		}).Info("AggregateDistance")
	}(time.Now())
	err = m.next.AggregateDistance(distance)
	return err
}

func (m *LogMiddleware) CalculateInvoice(obuid int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		var (
			distance float64
			amount   float64
		)
		if inv != nil {
			distance = inv.TotalDistance
			amount = inv.TotalAmount
		}
		logrus.WithFields(logrus.Fields{
			"took":     time.Since(start),
			"err":      err,
			"obuid":    obuid,
			"amount":   amount,
			"distance": distance,
		}).Info("CalculateInvoice")
	}(time.Now())
	inv, err = m.next.CalculateInvoice(obuid)
	return inv, err
}
