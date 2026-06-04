package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	serverName     = "hubitat-device-events"
	serverVersion  = "0.1.0"
	defaultBaseURL = "http://localhost:3000"
)

func main() {
	baseURL := os.Getenv("HUBITAT_API_BASE_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	client := NewHubitatClient(baseURL)

	s := server.NewMCPServer(
		serverName,
		serverVersion,
		server.WithToolCapabilities(true),
		server.WithInstructions(
			"Tools to read and manage smart-home devices and their events from a "+
				"HubitatDeviceEvents server. Device IDs are numeric strings. Event "+
				"time filters use ISO-8601 timestamps (e.g. 2024-01-01T00:00:00+00:00).",
		),
	)

	registerTools(s, client)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "%s: fatal: %v\n", serverName, err)
		os.Exit(1)
	}
}

func registerTools(s *server.MCPServer, client *HubitatClient) {
	s.AddTool(mcp.NewTool("list_devices",
		mcp.WithDescription("List all registered Hubitat devices (id, name, label, type, room)."),
		mcp.WithReadOnlyHintAnnotation(true),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		body, err := client.ListDevices()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(body), nil
	})

	s.AddTool(mcp.NewTool("get_device",
		mcp.WithDescription("Get a single device by its numeric device ID."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("id", mcp.Required(), mcp.Description("The numeric device ID.")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := req.RequireString("id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		body, err := client.GetDevice(id)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(body), nil
	})

	s.AddTool(mcp.NewTool("add_device",
		mcp.WithDescription("Register a new device. Provide the numeric ID plus metadata."),
		mcp.WithString("id", mcp.Required(), mcp.Description("The numeric device ID.")),
		mcp.WithString("name", mcp.Required(), mcp.Description("Device name (driver/type name).")),
		mcp.WithString("label", mcp.Required(), mcp.Description("Human-friendly device label.")),
		mcp.WithString("type", mcp.Required(), mcp.Description("Device type.")),
		mcp.WithString("room", mcp.Description("Room the device belongs to (optional).")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		d, err := deviceFromRequest(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		body, err := client.AddDevice(d)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(body), nil
	})

	s.AddTool(mcp.NewTool("update_device",
		mcp.WithDescription("Update an existing device's metadata. The device is matched by its ID."),
		mcp.WithString("id", mcp.Required(), mcp.Description("The numeric device ID.")),
		mcp.WithString("name", mcp.Required(), mcp.Description("Device name (driver/type name).")),
		mcp.WithString("label", mcp.Required(), mcp.Description("Human-friendly device label.")),
		mcp.WithString("type", mcp.Required(), mcp.Description("Device type.")),
		mcp.WithString("room", mcp.Description("Room the device belongs to (optional).")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		d, err := deviceFromRequest(req)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		body, err := client.UpdateDevice(d)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(body), nil
	})

	s.AddTool(mcp.NewTool("delete_device",
		mcp.WithDescription("Delete a device by ID. Fails if the device still has associated events."),
		mcp.WithString("id", mcp.Required(), mcp.Description("The numeric device ID.")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := req.RequireString("id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		body, err := client.DeleteDevice(id)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(body), nil
	})

	s.AddTool(mcp.NewTool("list_device_events",
		mcp.WithDescription("List all recorded device events across every device."),
		mcp.WithReadOnlyHintAnnotation(true),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		body, err := client.ListDeviceEvents()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(body), nil
	})

	s.AddTool(mcp.NewTool("get_device_events",
		mcp.WithDescription("Get events for a specific device, optionally filtered by an ISO-8601 start/end time range."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithString("id", mcp.Required(), mcp.Description("The numeric device ID.")),
		mcp.WithString("start", mcp.Description("ISO-8601 start timestamp, e.g. 2024-01-01T00:00:00+00:00.")),
		mcp.WithString("end", mcp.Description("ISO-8601 end timestamp, e.g. 2024-01-01T23:59:59+00:00.")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := req.RequireString("id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		start := req.GetString("start", "")
		end := req.GetString("end", "")
		if (start == "") != (end == "") {
			return mcp.NewToolResultError("both 'start' and 'end' must be provided together, or neither"), nil
		}
		body, err := client.GetDeviceEvents(id, start, end)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(body), nil
	})
}

func deviceFromRequest(req mcp.CallToolRequest) (Device, error) {
	id, err := req.RequireString("id")
	if err != nil {
		return Device{}, err
	}
	name, err := req.RequireString("name")
	if err != nil {
		return Device{}, err
	}
	label, err := req.RequireString("label")
	if err != nil {
		return Device{}, err
	}
	typ, err := req.RequireString("type")
	if err != nil {
		return Device{}, err
	}
	return Device{
		Id:    id,
		Name:  name,
		Label: label,
		Type:  typ,
		Room:  req.GetString("room", ""),
	}, nil
}
