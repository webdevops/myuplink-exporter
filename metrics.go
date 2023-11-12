package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

type (
	MyUplinkMetrics struct {
		system                *prometheus.GaugeVec
		systemDevice          *prometheus.GaugeVec
		systemDevicePoint     *prometheus.GaugeVec
		systemDevicePointEnum *prometheus.GaugeVec
	}
)

func NewMyUplinkMetrics(registry *prometheus.Registry) *MyUplinkMetrics {
	metrics := &MyUplinkMetrics{}

	metrics.system = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "myuplink_system_info",
			Help: "myUplink system information",
		},
		[]string{"systemID", "systemName", "country"},
	)
	registry.MustRegister(metrics.system)

	metrics.systemDevice = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "myuplink_system_device_info",
			Help: "myUplink system device information",
		},
		[]string{"systemID", "deviceID", "deviceName", "serialNumber", "connectionState", "firmwareVersion"},
	)
	registry.MustRegister(metrics.systemDevice)

	metrics.systemDevicePoint = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "myuplink_system_device_point",
			Help: "myUplink device metric point",
		},
		[]string{"systemID", "deviceID", "category", "parameterID", "parameterName", "parameterUnit"},
	)
	registry.MustRegister(metrics.systemDevicePoint)

	metrics.systemDevicePointEnum = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "myuplink_system_device_point_enum",
			Help: "myUplink device metric point enum value",
		},
		[]string{"systemID", "deviceID", "category", "parameterID", "parameterName", "parameterUnit", "valueText"},
	)
	registry.MustRegister(metrics.systemDevicePointEnum)

	return metrics
}
