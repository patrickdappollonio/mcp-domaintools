package resolver

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	resp "github.com/patrickdappollonio/mcp-domaintools/internal/response"
	"github.com/patrickdappollonio/mcp-domaintools/internal/utils"
)

// Config holds resolver configuration.
type Config struct {
	Timeout time.Duration
}

// resolverParams represents the parameters for hostname resolution.
type resolverParams struct {
	Hostname  string `json:"hostname"`
	IPVersion string `json:"ip_version"`
}

// HandleHostnameResolution resolves a hostname to its IP addresses.
func HandleHostnameResolution(ctx context.Context, request mcp.CallToolRequest, config *Config) (*mcp.CallToolResult, error) {
	var params resolverParams
	if err := request.BindArguments(&params); err != nil {
		return nil, fmt.Errorf("failed to parse tool input: %w", utils.ParseJSONUnmarshalError(err))
	}

	// Validate required parameters
	if params.Hostname == "" {
		return nil, fmt.Errorf("parameter \"hostname\" is required")
	}

	// Set default IP version if not provided
	if params.IPVersion == "" {
		params.IPVersion = "ipv4"
	}

	// Create context with timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()

	// Initialize response maps
	responseData := map[string]interface{}{
		"hostname":   params.Hostname,
		"timestamp":  time.Now().Format(time.RFC3339),
		"ip_version": params.IPVersion,
	}

	// Resolve based on IP version
	switch params.IPVersion {
	case "ipv4", "ipv6":
		addresses, err := lookupIPAddresses(ctxWithTimeout, params.Hostname, params.IPVersion)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve %s addresses: %w", params.IPVersion, err)
		}
		responseData[params.IPVersion+"_addresses"] = addresses

	default: // "both"
		// Get IPv4 addresses
		ipv4Addresses, err := lookupIPAddresses(ctxWithTimeout, params.Hostname, "ipv4")
		if err == nil {
			responseData["ipv4_addresses"] = ipv4Addresses
		} else {
			responseData["ipv4_error"] = err.Error()
		}

		// Get IPv6 addresses
		ipv6Addresses, err := lookupIPAddresses(ctxWithTimeout, params.Hostname, "ipv6")
		if err == nil {
			responseData["ipv6_addresses"] = ipv6Addresses
		} else {
			responseData["ipv6_error"] = err.Error()
		}
	}

	// Use the response package to handle JSON encoding and MCP tool result creation
	return resp.JSON(responseData)
}

// lookupIPAddresses handles the IP lookup for a specific IP version (IPv4 or IPv6).
func lookupIPAddresses(ctx context.Context, hostname, ipVersion string) ([]string, error) {
	networkType := "ip4"
	if ipVersion == "ipv6" {
		networkType = "ip6"
	}

	addrs, err := net.DefaultResolver.LookupIP(ctx, networkType, hostname)
	if err != nil {
		return nil, err
	}

	result := make([]string, len(addrs))
	for i, addr := range addrs {
		result[i] = addr.String()
	}

	return result, nil
}
