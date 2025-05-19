package server

import (
	"context"

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

	// Add handlers for the tools
	s.AddTool(localQueryTool, localDNSHandler)
	s.AddTool(remoteQueryTool, remoteDNSHandler)
	s.AddTool(whoisQueryTool, whoisHandler)

	return s, nil
}
