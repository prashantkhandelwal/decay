package middleware

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (

	// Upload metrics

	TotalFileUploadRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "total_file_uploaded_requests",
			Help: "Total number of file upload requests.",
		},
		[]string{"count"},
	)

	SuccessfulFileUploadRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "successful_file_upload_requests",
			Help: "Number of successful file upload requests.",
		},
		[]string{"count"},
	)

	FailedFileUploadRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "failed_file_upload_requests",
			Help: "Number of failed file upload requests.",
		},
		[]string{"count"},
	)

	// Server errors

	HttpRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

	HttpBadRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_bad_requests_total",
			Help: "Total number of bad HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

	HttpHeartbeatRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_heartbeat_requests_total",
			Help: "Total number of heartbeat HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)
)
