# HubitatDeviceEvents MCP Server

A [Model Context Protocol](https://modelcontextprotocol.io) (MCP) server that lets
an AI client read and manage smart-home devices and their events from a running
[HubitatDeviceEvents](../README) instance.

It is a **decoupled client**: it talks to the existing `/api/v1/` REST API over
HTTP, so the main Go app and this MCP server run as independent processes.

## Tools

| Tool | Description | REST call |
|---|---|---|
| `list_devices` | List all registered devices | `GET /api/v1/device` |
| `get_device` | Get one device by ID | `GET /api/v1/device/{id}` |
| `add_device` | Register a new device | `POST /api/v1/device` |
| `update_device` | Update a device's metadata | `PUT /api/v1/device` |
| `delete_device` | Delete a device by ID | `DELETE /api/v1/device/{id}` |
| `list_device_events` | List all events across all devices | `GET /api/v1/device/event` |
| `get_device_events` | Events for one device, optional `start`/`end` ISO range | `GET /api/v1/device/{id}/event` |

## Configuration

| Env var | Purpose | Default |
|---|---|---|
| `HUBITAT_API_BASE_URL` | Base URL of the HubitatDeviceEvents REST API | `http://localhost:3000` |

## Build

```sh
cd mcp
go build -o hubitat-mcp.exe .   # Windows
# go build -o hubitat-mcp .     # Linux/macOS
```

This produces a self-contained `hubitat-mcp` binary that speaks MCP over
**stdio** — the transport used by most local AI clients.

## Attaching to a local AI client

Make sure your HubitatDeviceEvents server is running and reachable (default
`http://localhost:3000`), then point your client at the built binary.

### VS Code (`.vscode/mcp.json`) or GitHub Copilot

```json
{
  "servers": {
    "hubitat-device-events": {
      "type": "stdio",
      "command": "C:\\Local-Projects\\HubitatDeviceEvents\\mcp\\hubitat-mcp.exe",
      "env": {
        "HUBITAT_API_BASE_URL": "http://localhost:3000"
      }
    }
  }
}
```

### Claude Desktop (`claude_desktop_config.json`)

```json
{
  "mcpServers": {
    "hubitat-device-events": {
      "command": "C:\\Local-Projects\\HubitatDeviceEvents\\mcp\\hubitat-mcp.exe",
      "env": {
        "HUBITAT_API_BASE_URL": "http://localhost:3000"
      }
    }
  }
}
```

Use the absolute path to the binary. On Linux/macOS drop the `.exe` and use a
POSIX path. After saving the config, restart the client; the seven tools above
should appear.

## Quick manual test (stdio)

```sh
printf '%s\n' \
 '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"t","version":"1"}}}' \
 '{"jsonrpc":"2.0","method":"notifications/initialized"}' \
 '{"jsonrpc":"2.0","id":2,"method":"tools/list"}' | ./hubitat-mcp
```

You should see a JSON-RPC response listing the available tools.

## Notes

- Device IDs are numeric strings (e.g. `"42"`); the device JSON uses an
  uppercase `Id` field to match the Hubitat API.
- Event time filters use ISO-8601 timestamps, e.g.
  `2024-01-01T00:00:00+00:00`. Provide `start` and `end` together or not at all.
- This is a separate Go module (`mcp/go.mod`) so its dependencies stay isolated
  from the main application.
