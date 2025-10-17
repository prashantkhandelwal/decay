package middleware

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	HttpRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

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
)
