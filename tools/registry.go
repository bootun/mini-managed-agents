package tools

import (
	"context"
	"fmt"
)

type Handler func(ctx context.Context, args map[string]any) (string, error)

func GetTools() []map[string]any {
	return []map[string]any{
		weatherAlertsTool(),
		getLocationTool(),
		getIPAddressTool(),
	}
}

func GetHandler(name string) (Handler, error) {
	switch name {
	case "get_weather_alerts":
		return GetWeatherAlerts, nil
	case "get_location_info":
		return GetLocationInfo, nil
	case "get_ip_address":
		return GetIPAddress, nil
	default:
		return nil, fmt.Errorf("unknown tool name: %s", name)
	}
}
