package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/webdevops/myuplink-exporter/myuplink"
)

const (
	DefaultTimeout = 30
)

func myuplinkProbe(w http.ResponseWriter, r *http.Request) {
	var err error
	var timeoutSeconds float64

	// startTime := time.Now()
	contextLogger := buildContextLoggerFromRequest(r)
	registry := prometheus.NewRegistry()

	// If a timeout is configured via the Prometheus header, add it to the request.
	timeoutSeconds, err = getPrometheusTimeout(r, DefaultTimeout)
	if err != nil {
		contextLogger.Error(err.Error())
		http.Error(w, fmt.Sprintf("failed to parse timeout from Prometheus header: %s", err), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds*float64(time.Second)))
	defer cancel()
	r = r.WithContext(ctx)

	// use timeout as max cache time as mostly it's also the scrape time
	cacheTime := time.Duration(timeoutSeconds) * time.Second

	if v := r.URL.Query().Get("cache"); v != "" {
		cacheTime, err = time.ParseDuration(v)
		if err != nil {
			contextLogger.Error(err.Error())
			http.Error(w, fmt.Sprintf("failed to parse cache from query param: %s", err), http.StatusBadRequest)
			return
		}
	}

	metricSystem := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "myuplink_system_info",
			Help: "",
		},
		[]string{"systemID", "systemName", "country"},
	)
	registry.MustRegister(metricSystem)

	metricSystemDevice := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "myuplink_system_device_info",
			Help: "",
		},
		[]string{"systemID", "deviceID", "deviceName", "serialNumber", "connectionState", "firmwareVersion"},
	)
	registry.MustRegister(metricSystemDevice)

	metricSystemDevicePoint := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "myuplink_system_device_point",
			Help: "",
		},
		[]string{"systemID", "deviceID", "category", "parameterID", "parameterName", "parameterUnit"},
	)
	registry.MustRegister(metricSystemDevicePoint)

	systemList, err := cacheResult(
		"systems",
		func() (interface{}, error) {
			return myuplinkClient.GetSystems(ctx)
		},
	)
	if err != nil {
		contextLogger.Error(err.Error())
		http.Error(w, fmt.Sprintf("failed to fetch system list from myUplink: %s", err), http.StatusBadRequest)
		return
	}

	for _, system := range systemList.(*myuplink.ResultSystems).Systems {
		metricSystem.With(prometheus.Labels{
			"systemID":   system.SystemID,
			"systemName": system.Name,
			"country":    system.Country,
		}).Set(1)

		for _, device := range system.Devices {
			metricSystemDevice.With(prometheus.Labels{
				"systemID":        system.SystemID,
				"deviceID":        device.ID,
				"deviceName":      device.Product.Name,
				"serialNumber":    device.Product.SerialNumber,
				"connectionState": device.ConnectionState,
				"firmwareVersion": device.CurrentFwVersion,
			}).Set(1)

			devicePoints, err := cacheResultWithDuration(
				fmt.Sprintf("devicePoints:%s", device.ID),
				cacheTime,
				func() (interface{}, error) {
					return myuplinkClient.GetSystemDevicePoints(ctx, device.ID)
				},
			)
			if err != nil {
				contextLogger.Error(err.Error())
				http.Error(w, fmt.Sprintf("failed to fetch device points from myUplink: %s", err), http.StatusBadRequest)
				return
			}

			for _, devicePoint := range *devicePoints.(*myuplink.SystemDevicePoints) {
				metricSystemDevicePoint.With(prometheus.Labels{
					"systemID":      system.SystemID,
					"deviceID":      device.ID,
					"category":      devicePoint.Category,
					"parameterID":   devicePoint.ParameterID,
					"parameterName": devicePoint.ParameterName,
					"parameterUnit": devicePoint.ParameterUnit,
				}).Set(devicePoint.Value)
			}
		}
	}

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func buildContextLoggerFromRequest(r *http.Request) *zap.SugaredLogger {
	return logger.With(zap.String("requestPath", r.URL.Path))
}

func getPrometheusTimeout(r *http.Request, defaultTimeout float64) (timeout float64, err error) {
	// If a timeout is configured via the Prometheus header, add it to the request.
	if v := r.Header.Get("X-Prometheus-Scrape-Timeout-Seconds"); v != "" {
		timeout, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return
		}
	}
	if timeout == 0 {
		timeout = defaultTimeout
	}

	return
}

// cacheResult caches template function results (eg. Azure REST API resource information)
func cacheResult(cacheKey string, callback func() (interface{}, error)) (interface{}, error) {
	if val, ok := globalCache.Get(cacheKey); ok {
		return val, nil
	}

	ret, err := callback()
	if err != nil {
		return nil, err
	}

	globalCache.SetDefault(cacheKey, ret)

	return ret, nil
}

func cacheResultWithDuration(cacheKey string, cacheTime time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	if val, ok := globalCache.Get(cacheKey); ok {
		return val, nil
	}

	ret, err := callback()
	if err != nil {
		return nil, err
	}

	globalCache.Set(cacheKey, ret, cacheTime)

	return ret, nil
}
