package whois

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/likexian/whois"
	"github.com/mark3labs/mcp-go/mcp"
)

// Config holds WHOIS configuration.
type Config struct {
	CustomServer string
}

// HandleWhoisQuery processes WHOIS queries.
func HandleWhoisQuery(ctx context.Context, request mcp.CallToolRequest, config *Config) (*mcp.CallToolResult, error) {
	domain := mcp.ParseString(request, "domain", "")
	if domain == "" {
		return nil, fmt.Errorf("parameter \"domain\" is required")
	}

	// Clean and validate domain format
	domain = strings.TrimSpace(domain)
	if strings.Contains(domain, "..") || strings.HasPrefix(domain, ".") {
		return nil, fmt.Errorf("invalid domain format: %q", domain)
	}

	var result string
	var err error

	// Use custom server if provided, otherwise use default
	if config.CustomServer != "" {
		result, err = whois.Whois(domain, config.CustomServer)
	} else {
		result, err = whois.Whois(domain)
	}

	if err != nil {
		return nil, fmt.Errorf("WHOIS query failed: %s", err)
	}

	// Format response as JSON
	response := map[string]interface{}{
		"domain": domain,
		"result": result,
	}

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("error generating JSON: %s", err)
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}
