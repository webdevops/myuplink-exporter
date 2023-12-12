package myuplink

import (
	"strings"
	"time"
)

type (
	ResultSystems struct {
		Page         int `json:"page"`
		ItemsPerPage int `json:"itemsPerPage"`
		NumItems     int `json:"numItems"`
		Systems      []struct {
			SystemID      string         `json:"systemId"`
			Name          string         `json:"name"`
			SecurityLevel string         `json:"securityLevel"`
			HasAlarm      bool           `json:"hasAlarm"`
			Country       string         `json:"country"`
			Devices       []SystemDevice `json:"devices"`
		} `json:"systems"`
	}

	SystemDevice struct {
		ID               string `json:"id"`
		ConnectionState  string `json:"connectionState"`
		CurrentFwVersion string `json:"currentFwVersion"`
		Product          struct {
			SerialNumber string `json:"serialNumber"`
			Name         string `json:"name"`
		} `json:"product"`
	}

	SystemDeviceFirmwareInfo struct {
		DeviceID         string `json:"deviceId"`
		FirmwareID       string `json:"firmwareId"`
		CurrentFwVersion string `json:"currentFwVersion"`
		PendingFwVersion string `json:"pendingFwVersion"`
		DesiredFwVersion string `json:"desiredFwVersion"`
	}

	SystemDevicePoints []struct {
		Category            string    `json:"category"`
		ParameterID         string    `json:"parameterId"`
		ParameterName       string    `json:"parameterName"`
		ParameterUnit       string    `json:"parameterUnit"`
		Writable            bool      `json:"writable"`
		Timestamp           time.Time `json:"timestamp"`
		Value               *float64  `json:"value"`
		StrVal              string    `json:"strVal"`
		SmartHomeCategories []string  `json:"smartHomeCategories"`
		MinValue            float64   `json:"minValue"`
		MaxValue            float64   `json:"maxValue"`
		StepValue           float64   `json:"stepValue"`
		EnumValues          []struct {
			Value string `json:"value"`
			Text  string `json:"text"`
			Icon  string `json:"icon"`
		} `json:"enumValues"`
		ScaleValue string `json:"scaleValue"`
		ZoneID     string `json:"zoneId"`
	}
)

func (d *SystemDevice) IsConnectionStateAllowed(allowedValues []string) bool {
	for _, allowedVal := range allowedValues {
		if strings.EqualFold(d.ConnectionState, allowedVal) {
			return true
		}
	}

	return false
}
