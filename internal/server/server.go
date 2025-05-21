package server

import (
	"context"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/patrickdappollonio/mcp-domaintools/internal/dns"
	"github.com/patrickdappollonio/mcp-domaintools/internal/resolver"
	"github.com/patrickdappollonio/mcp-domaintools/internal/whois"
)

// DomainToolsConfig contains configuration for the domain tools.
type DomainToolsConfig struct {
	QueryConfig     *dns.QueryConfig
	WhoisConfig     *whois.Config
	ResolverConfig  *resolver.Config
	Version         string
}

// SetupTools creates and configures the domain query tools.
func SetupTools(config *DomainToolsConfig) (*server.MCPServer, error) {
	// Create a new MCP server
	s := server.NewMCPServer(
		"DNS and WHOIS Query Tools",
		config.Version,
		server.WithRecovery(),
	)
	
	// Initialize resolver config if not provided
	if config.ResolverConfig == nil {
		config.ResolverConfig = &resolver.Config{
			Timeout: 5 * time.Second,
		}
	}

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
		return resolver.HandleHostnameResolution(ctx, request, config.ResolverConfig)
	}

	// Add handlers for the tools
	s.AddTool(localQueryTool, localDNSHandler)
	s.AddTool(remoteQueryTool, remoteDNSHandler)
	s.AddTool(whoisQueryTool, whoisHandler)
	s.AddTool(resolveHostTool, resolveHostHandler)

	return s, nil
}


