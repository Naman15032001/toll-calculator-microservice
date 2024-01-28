package main

import (
	"time"

	"github.com/Naman15032001/tolling/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next DataProducer
}

func NewLogMiddleware(next DataProducer) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}

func (l *LogMiddleware) ProduceData(data types.OBUDATA) error {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"obuid": data.OBUID,
			"lat":   data.Lat,
			"long":  data.Long,
			"took": time.Since(start),
		}).Info("producing to kafka")
	}(time.Now())
	return l.next.ProduceData(data)
}
