package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/patrickdappollonio/mcp-domaintools/internal/dns"
	internalServer "github.com/patrickdappollonio/mcp-domaintools/internal/server"
	"github.com/patrickdappollonio/mcp-domaintools/internal/whois"
)

var (
	remoteServerAddress string
	customWhoisServer   string
	timeout             time.Duration
	version             = "dev"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(1)
	}
}

func run() error {
	flag.StringVar(&remoteServerAddress, "remote-server-address", "", "Custom DNS-over-HTTPS server address")
	flag.StringVar(&customWhoisServer, "custom-whois-server", "", "Custom WHOIS server address")
	flag.DurationVar(&timeout, "timeout", 5*time.Second, "Timeout for DNS queries")

	flag.Parse()

	// Create DNS query configuration
	queryConfig := &dns.QueryConfig{
		Timeout:             timeout,
		RemoteServerAddress: remoteServerAddress,
	}

	// Create WHOIS configuration
	whoisConfig := &whois.Config{
		CustomServer: customWhoisServer,
	}

	// Setup domain tools
	s, err := internalServer.SetupTools(&internalServer.DomainToolsConfig{
		QueryConfig: queryConfig,
		WhoisConfig: whoisConfig,
		Version:     version,
	})
	if err != nil {
		return fmt.Errorf("error setting up domain tools: %w", err)
	}

	// Start the server
	if err := server.ServeStdio(s); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
