package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/webdevops/myuplink-exporter/myuplink"
)

type (
	TotalParameterCache struct {
		lock           sync.RWMutex
		deviceParamMap map[string]*TotalParameterMetric
	}

	TotalParameterMetric struct {
		updated time.Time
		value   float64
	}
)

var (
	totalParamCache TotalParameterCache
)

func (c *TotalParameterCache) Init() {
	c.deviceParamMap = map[string]*TotalParameterMetric{}
}

func (c *TotalParameterCache) getParameterValue(deviceID string, parameterID string, devicePoint myuplink.SystemDevicePoint) float64 {
	key := fmt.Sprintf(
		"%s:%s",
		deviceID,
		parameterID,
	)

	// exists check (eg. new parameter)
	if _, exists := c.deviceParamMap[key]; !exists {
		c.lock.Lock()
		c.deviceParamMap[key] = &TotalParameterMetric{
			updated: devicePoint.Timestamp,
			value:   *devicePoint.Value,
		}
		c.lock.Unlock()
	}

	// increase check
	if c.deviceParamMap[key].updated != devicePoint.Timestamp {
		c.lock.Lock()
		c.deviceParamMap[key].updated = devicePoint.Timestamp
		c.deviceParamMap[key].value += *devicePoint.Value
		c.lock.Unlock()
	}

	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.deviceParamMap[key].value
}
