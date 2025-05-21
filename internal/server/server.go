package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/patrickdappollonio/mcp-domaintools/internal/dns"
	"github.com/patrickdappollonio/mcp-domaintools/internal/whois"
)

// DomainToolsConfig contains configuration for the domain tools.
type DomainToolsConfig struct {
	QueryConfig *dns.QueryConfig
	WhoisConfig *whois.Config
	Version     string
}

// SetupTools creates and configures the domain query tools.
func SetupTools(config *DomainToolsConfig) (*server.MCPServer, error) {
	// Create a new MCP server
	s := server.NewMCPServer(
		"DNS and WHOIS Query Tools",
		config.Version,
		server.WithRecovery(),
	)

	// Add local DNS query tool
	localQueryTool := mcp.NewTool("local_dns_query",
		mcp.WithDescription("Perform DNS queries using local OS-defined DNS servers"),
		mcp.WithString("domain",
			mcp.Required(),
			mcp.Description("The domain name to query (e.g., example.com)"),
		),
		mcp.WithString("record_type",
			mcp.Required(),
			mcp.Description("The type of DNS record to query"),
			mcp.Enum("A", "AAAA", "CNAME", "MX", "NS", "PTR", "SOA", "SRV", "TXT"),
		),
	)

	// Add remote DNS query tool
	remoteQueryTool := mcp.NewTool("remote_dns_query",
		mcp.WithDescription("Perform DNS queries using remote DNS-over-HTTPS servers (Google and Cloudflare)"),
		mcp.WithString("domain",
			mcp.Required(),
			mcp.Description("The domain name to query (e.g., example.com)"),
		),
		mcp.WithString("record_type",
			mcp.Required(),
			mcp.Description("The type of DNS record to query"),
			mcp.Enum("A", "AAAA", "CNAME", "MX", "NS", "PTR", "SOA", "SRV", "TXT"),
		),
	)

	// Add WHOIS query tool
	whoisQueryTool := mcp.NewTool("whois_query",
		mcp.WithDescription("Perform WHOIS lookups to get domain registration information"),
		mcp.WithString("domain",
			mcp.Required(),
			mcp.Description("The domain name to query (e.g., example.com)"),
		),
	)

	// Add hostname to IP resolution tool
	resolveHostTool := mcp.NewTool("resolve_hostname",
		mcp.WithDescription("Convert a hostname to its corresponding IP addresses"),
		mcp.WithString("hostname",
			mcp.Required(),
			mcp.Description("The hostname to resolve (e.g., example.com)"),
		),
		mcp.WithString("ip_version",
			mcp.Description("IP version to resolve (ipv4, ipv6, or both); defaults to ipv4"),
			mcp.Enum("ipv4", "ipv6", "both"),
		),
	)

	// Create handler wrappers
	localDNSHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return dns.HandleLocalDNSQuery(ctx, request, config.QueryConfig)
	}

	remoteDNSHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return dns.HandleRemoteDNSQuery(ctx, request, config.QueryConfig)
	}

	whoisHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return whois.HandleWhoisQuery(ctx, request, config.WhoisConfig)
	}

	resolveHostHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleHostnameResolution(ctx, request, config.QueryConfig)
	}

	// Add handlers for the tools
	s.AddTool(localQueryTool, localDNSHandler)
	s.AddTool(remoteQueryTool, remoteDNSHandler)
	s.AddTool(whoisQueryTool, whoisHandler)
	s.AddTool(resolveHostTool, resolveHostHandler)

	return s, nil
}

// handleHostnameResolution resolves a hostname to its IP addresses.
func handleHostnameResolution(ctx context.Context, request mcp.CallToolRequest, config *dns.QueryConfig) (*mcp.CallToolResult, error) {
	hostname := mcp.ParseString(request, "hostname", "")
	if hostname == "" {
		return nil, fmt.Errorf("parameter \"hostname\" is required")
	}

	ipVersion := mcp.ParseString(request, "ip_version", "ipv4")

	// Create context with timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()

	// Initialize response maps
	response := map[string]interface{}{
		"hostname":  hostname,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	// Resolve based on IP version
	switch ipVersion {
	case "ipv4":
		ipv4Addrs, err := net.DefaultResolver.LookupIP(ctxWithTimeout, "ip4", hostname)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve IPv4 addresses: %w", err)
		}

		ipv4Strings := make([]string, len(ipv4Addrs))
		for i, addr := range ipv4Addrs {
			ipv4Strings[i] = addr.String()
		}

		response["ipv4_addresses"] = ipv4Strings
		response["ip_version"] = "ipv4"

	case "ipv6":
		ipv6Addrs, err := net.DefaultResolver.LookupIP(ctxWithTimeout, "ip6", hostname)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve IPv6 addresses: %w", err)
		}

		ipv6Strings := make([]string, len(ipv6Addrs))
		for i, addr := range ipv6Addrs {
			ipv6Strings[i] = addr.String()
		}

		response["ipv6_addresses"] = ipv6Strings
		response["ip_version"] = "ipv6"

	default: // "both"
		// Get IPv4 addresses
		ipv4Addrs, err := net.DefaultResolver.LookupIP(ctxWithTimeout, "ip4", hostname)
		if err == nil {
			ipv4Strings := make([]string, len(ipv4Addrs))
			for i, addr := range ipv4Addrs {
				ipv4Strings[i] = addr.String()
			}
			response["ipv4_addresses"] = ipv4Strings
		} else {
			response["ipv4_error"] = err.Error()
		}

		// Get IPv6 addresses
		ipv6Addrs, err := net.DefaultResolver.LookupIP(ctxWithTimeout, "ip6", hostname)
		if err == nil {
			ipv6Strings := make([]string, len(ipv6Addrs))
			for i, addr := range ipv6Addrs {
				ipv6Strings[i] = addr.String()
			}
			response["ipv6_addresses"] = ipv6Strings
		} else {
			response["ipv6_error"] = err.Error()
		}

		response["ip_version"] = "both"
	}

	// Marshal the response to JSON
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("error generating JSON: %w", err)
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}
