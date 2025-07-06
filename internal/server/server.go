package server

import (
	"context"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/patrickdappollonio/mcp-domaintools/internal/dns"
	"github.com/patrickdappollonio/mcp-domaintools/internal/http_ping"
	"github.com/patrickdappollonio/mcp-domaintools/internal/ping"
	"github.com/patrickdappollonio/mcp-domaintools/internal/resolver"
	"github.com/patrickdappollonio/mcp-domaintools/internal/tls"
	"github.com/patrickdappollonio/mcp-domaintools/internal/whois"
)

// DomainToolsConfig contains configuration for the domain tools.
type DomainToolsConfig struct {
	QueryConfig    *dns.QueryConfig
	WhoisConfig    *whois.Config
	ResolverConfig *resolver.Config
	PingConfig     *ping.Config
	HTTPPingConfig *http_ping.Config
	TLSConfig      *tls.Config
	Version        string
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

	// Initialize ping config if not provided
	if config.PingConfig == nil {
		config.PingConfig = &ping.Config{
			Timeout: 5 * time.Second,
			Count:   4,
		}
	}

	// Initialize HTTP ping config if not provided
	if config.HTTPPingConfig == nil {
		config.HTTPPingConfig = &http_ping.Config{
			Timeout: 10 * time.Second,
			Count:   1,
		}
	}

	// Initialize TLS config if not provided
	if config.TLSConfig == nil {
		config.TLSConfig = &tls.Config{
			Timeout: 10 * time.Second,
			Port:    443,
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
			mcp.Description("The type of DNS record to query; defaults to A"),
			mcp.Enum("A", "AAAA", "CNAME", "MX", "NS", "PTR", "SOA", "SRV", "TXT"),
			mcp.DefaultString("A"),
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
			mcp.Description("The type of DNS record to query; defaults to A"),
			mcp.Enum("A", "AAAA", "CNAME", "MX", "NS", "PTR", "SOA", "SRV", "TXT"),
			mcp.DefaultString("A"),
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
			mcp.DefaultString("ipv4"),
		),
	)

	// Add ping tool
	pingTool := mcp.NewTool("ping",
		mcp.WithDescription("Perform ping operations to test connectivity and measure response times to a host"),
		mcp.WithString("target",
			mcp.Required(),
			mcp.Description("The hostname or IP address to ping (e.g., example.com or 8.8.8.8)"),
		),
		mcp.WithNumber("count",
			mcp.Description("Number of ping packets to send; defaults to 4"),
			mcp.DefaultNumber(4),
		),
	)

	// Add TLS certificate check tool
	tlsCheckTool := mcp.NewTool("tls_certificate_check",
		mcp.WithDescription("Check TLS certificate chain for a domain to analyze certificate validity, expiration, and chain structure"),
		mcp.WithString("domain",
			mcp.Required(),
			mcp.Description("The domain name to check TLS certificate for (e.g., example.com)"),
		),
		mcp.WithNumber("port",
			mcp.Description("Port to connect to for TLS check; defaults to 443"),
			mcp.DefaultNumber(443),
		),
		mcp.WithBoolean("include_chain",
			mcp.Description("Whether to include the full certificate chain in the response; defaults to true"),
		),
		mcp.WithBoolean("check_expiry",
			mcp.Description("Whether to check certificate expiration and provide warnings; defaults to true"),
		),
		mcp.WithString("server_name",
			mcp.Description("Server name for SNI (Server Name Indication); defaults to the domain name"),
		),
	)

	// Add HTTP ping tool
	httpPingTool := mcp.NewTool("http_ping",
		mcp.WithDescription("Perform HTTP ping operations to test connectivity and measure response times to HTTP endpoints"),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("The URL to ping (e.g., https://api.example.com/users)"),
		),
		mcp.WithString("method",
			mcp.Description("HTTP method to use; defaults to GET"),
			mcp.Enum("GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "PATCH"),
			mcp.DefaultString("GET"),
		),
		mcp.WithNumber("count",
			mcp.Description("Number of HTTP requests to send; defaults to 1"),
			mcp.DefaultNumber(1),
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

	pingHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return ping.HandlePing(ctx, request, config.PingConfig)
	}

	tlsCheckHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tls.HandleTLSCheck(ctx, request, config.TLSConfig)
	}

	httpPingHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return http_ping.HandleHTTPPing(ctx, request, config.HTTPPingConfig)
	}

	// Add handlers for the tools
	s.AddTool(localQueryTool, localDNSHandler)
	s.AddTool(remoteQueryTool, remoteDNSHandler)
	s.AddTool(whoisQueryTool, whoisHandler)
	s.AddTool(resolveHostTool, resolveHostHandler)
	s.AddTool(pingTool, pingHandler)
	s.AddTool(tlsCheckTool, tlsCheckHandler)
	s.AddTool(httpPingTool, httpPingHandler)

	return s, nil
}
