package handlers

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strconv"
	"strings"
)

var (
	putMetrics = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "cache_put",
		Help: "Total number of saved items",
	}, []string{"path0", "path1", "code"})
	getMetrics = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "cache_get",
		Help: "Total number of requested items",
	}, []string{"path0", "path1", "code"})
	delMetrics = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "cache_del",
		Help: "Total number of deleted items",
	}, []string{"path0", "path1", "code"})
)

func IncCounter(path string, code int, method string) {
	var counter *prometheus.CounterVec
	switch method {
	case "GET":
		counter = getMetrics
		break
	case "PUT":
		counter = putMetrics
		break
	case "DELETE":
		counter = delMetrics
		break
	default:
		return
	}
	counter.With(getPrometheusLabels(path, code)).Inc()
}

func getPrometheusLabels(path string, code int) prometheus.Labels {
	labels := prometheus.Labels{"path0": "", "path1": "", "code": strconv.Itoa(code)}
	if path == "" {
		return labels
	}
	parts := strings.Split(path, "/")
	if len(parts) > 2 {
		labels["path0"] = parts[1]
	}
	if len(parts) > 3 {
		labels["path1"] = parts[2]
	}
	return labels
}
