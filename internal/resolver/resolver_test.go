package resolver

import (
	"context"
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleHostnameResolution_ParameterValidation(t *testing.T) {
	ctx := context.Background()
	config := &Config{Timeout: 5 * time.Second}

	t.Run("missing hostname parameter", func(t *testing.T) {
		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Arguments: map[string]interface{}{
					"ip_version": "ipv4",
				},
			},
		}

		result, err := HandleHostnameResolution(ctx, request, config)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), `parameter "hostname" is required`)
	})

	t.Run("empty hostname parameter", func(t *testing.T) {
		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Arguments: map[string]interface{}{
					"hostname":   "",
					"ip_version": "ipv4",
				},
			},
		}

		result, err := HandleHostnameResolution(ctx, request, config)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), `parameter "hostname" is required`)
	})

	t.Run("invalid json arguments", func(t *testing.T) {
		// This would simulate malformed JSON arguments
		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Arguments: map[string]interface{}{
					"hostname": 12345, // Wrong type
				},
			},
		}

		result, err := HandleHostnameResolution(ctx, request, config)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to parse tool input")
	})
}

func TestHandleHostnameResolution_DefaultIPVersion(t *testing.T) {
	ctx := context.Background()
	config := &Config{Timeout: 5 * time.Second}

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"hostname": "localhost",
				// ip_version is omitted - should default to ipv4
			},
		},
	}

	result, err := HandleHostnameResolution(ctx, request, config)

	require.NoError(t, err)
	require.NotNil(t, result)

	// Parse the JSON response
	require.Len(t, result.Content, 1)
	textContent, ok := mcp.AsTextContent(result.Content[0])
	require.True(t, ok)

	var responseData map[string]interface{}
	err = json.Unmarshal([]byte(textContent.Text), &responseData)
	require.NoError(t, err)

	assert.Equal(t, "localhost", responseData["hostname"])
	assert.Equal(t, "ipv4", responseData["ip_version"])
	assert.Equal(t, false, responseData["failed"])
}

func TestHandleHostnameResolution_SuccessfulIPv4Resolution(t *testing.T) {
	ctx := context.Background()
	config := &Config{Timeout: 5 * time.Second}

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"hostname":   "localhost",
				"ip_version": "ipv4",
			},
		},
	}

	result, err := HandleHostnameResolution(ctx, request, config)

	require.NoError(t, err)
	require.NotNil(t, result)

	// Parse the JSON response
	require.Len(t, result.Content, 1)
	textContent, ok := mcp.AsTextContent(result.Content[0])
	require.True(t, ok)

	var responseData map[string]interface{}
	err = json.Unmarshal([]byte(textContent.Text), &responseData)
	require.NoError(t, err)

	assert.Equal(t, "localhost", responseData["hostname"])
	assert.Equal(t, "ipv4", responseData["ip_version"])
	assert.Equal(t, false, responseData["failed"])
	assert.Contains(t, responseData, "ipv4_addresses")
	assert.NotContains(t, responseData, "error")

	// Check that we got IPv4 addresses
	addresses, ok := responseData["ipv4_addresses"].([]interface{})
	assert.True(t, ok)
	assert.NotEmpty(t, addresses)
}

func TestHandleHostnameResolution_SuccessfulIPv6Resolution(t *testing.T) {
	ctx := context.Background()
	config := &Config{Timeout: 5 * time.Second}

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"hostname":   "localhost",
				"ip_version": "ipv6",
			},
		},
	}

	result, err := HandleHostnameResolution(ctx, request, config)

	require.NoError(t, err)
	require.NotNil(t, result)

	// Parse the JSON response
	require.Len(t, result.Content, 1)
	textContent, ok := mcp.AsTextContent(result.Content[0])
	require.True(t, ok)

	var responseData map[string]interface{}
	err = json.Unmarshal([]byte(textContent.Text), &responseData)
	require.NoError(t, err)

	assert.Equal(t, "localhost", responseData["hostname"])
	assert.Equal(t, "ipv6", responseData["ip_version"])
	assert.Equal(t, false, responseData["failed"])

	// IPv6 might not always be available, so we check if it's there or if there's an error
	if ipv6Addresses, exists := responseData["ipv6_addresses"]; exists {
		addresses, ok := ipv6Addresses.([]interface{})
		assert.True(t, ok)
		assert.NotEmpty(t, addresses)
		assert.NotContains(t, responseData, "error")
	} else {
		// If IPv6 is not available, there should be an error
		assert.Contains(t, responseData, "error")
		assert.Equal(t, true, responseData["failed"])
	}
}

func TestHandleHostnameResolution_BothVersions(t *testing.T) {
	ctx := context.Background()
	config := &Config{Timeout: 5 * time.Second}

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"hostname":   "localhost",
				"ip_version": "both",
			},
		},
	}

	result, err := HandleHostnameResolution(ctx, request, config)

	require.NoError(t, err)
	require.NotNil(t, result)

	// Parse the JSON response
	require.Len(t, result.Content, 1)
	textContent, ok := mcp.AsTextContent(result.Content[0])
	require.True(t, ok)

	var responseData map[string]interface{}
	err = json.Unmarshal([]byte(textContent.Text), &responseData)
	require.NoError(t, err)

	assert.Equal(t, "localhost", responseData["hostname"])
	assert.Equal(t, "both", responseData["ip_version"])

	// Should have IPv4 addresses (localhost should always resolve to IPv4)
	assert.Contains(t, responseData, "ipv4_addresses")
	ipv4Addresses, ok := responseData["ipv4_addresses"].([]interface{})
	assert.True(t, ok)
	assert.NotEmpty(t, ipv4Addresses)
}

func TestHandleHostnameResolution_NonExistentDomain(t *testing.T) {
	ctx := context.Background()
	config := &Config{Timeout: 5 * time.Second}

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"hostname":   "this-domain-definitely-does-not-exist-12345.invalid",
				"ip_version": "ipv4",
			},
		},
	}

	result, err := HandleHostnameResolution(ctx, request, config)

	// Should NOT return an error - should return JSON with failed: true
	require.NoError(t, err)
	require.NotNil(t, result)

	// Parse the JSON response
	require.Len(t, result.Content, 1)
	textContent, ok := mcp.AsTextContent(result.Content[0])
	require.True(t, ok)

	var responseData map[string]interface{}
	err = json.Unmarshal([]byte(textContent.Text), &responseData)
	require.NoError(t, err)

	assert.Equal(t, "this-domain-definitely-does-not-exist-12345.invalid", responseData["hostname"])
	assert.Equal(t, "ipv4", responseData["ip_version"])
	assert.Equal(t, true, responseData["failed"])
	assert.Contains(t, responseData, "error")
	assert.NotContains(t, responseData, "ipv4_addresses")

	errorMsg, ok := responseData["error"].(string)
	assert.True(t, ok)
	assert.NotEmpty(t, errorMsg)
}

func TestHandleHostnameResolution_NonExistentDomainBothVersions(t *testing.T) {
	ctx := context.Background()
	config := &Config{Timeout: 5 * time.Second}

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"hostname":   "this-domain-definitely-does-not-exist-67890.invalid",
				"ip_version": "both",
			},
		},
	}

	result, err := HandleHostnameResolution(ctx, request, config)

	// Should NOT return an error - should return JSON with failed: true
	require.NoError(t, err)
	require.NotNil(t, result)

	// Parse the JSON response
	require.Len(t, result.Content, 1)
	textContent, ok := mcp.AsTextContent(result.Content[0])
	require.True(t, ok)

	var responseData map[string]interface{}
	err = json.Unmarshal([]byte(textContent.Text), &responseData)
	require.NoError(t, err)

	assert.Equal(t, "this-domain-definitely-does-not-exist-67890.invalid", responseData["hostname"])
	assert.Equal(t, "both", responseData["ip_version"])
	assert.Equal(t, true, responseData["failed"])
	assert.Contains(t, responseData, "error")
	assert.NotContains(t, responseData, "ipv4_addresses")
	assert.NotContains(t, responseData, "ipv6_addresses")
}

func TestHandleHostnameResolution_ContextTimeout(t *testing.T) {
	// Create a context that times out very quickly
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	config := &Config{Timeout: 1 * time.Nanosecond}

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"hostname":   "google.com",
				"ip_version": "ipv4",
			},
		},
	}

	// This should likely timeout, but the exact behavior might vary by system
	result, err := HandleHostnameResolution(ctx, request, config)

	if err != nil {
		// If it returns an error, it should be a timeout/context error (critical error)
		assert.Contains(t, err.Error(), "failed to resolve ipv4 addresses")
	} else {
		// If it returns a result, it might have succeeded very quickly or failed with host not found
		require.NotNil(t, result)

		require.Len(t, result.Content, 1)
		textContent, ok := mcp.AsTextContent(result.Content[0])
		require.True(t, ok)

		var responseData map[string]interface{}
		err = json.Unmarshal([]byte(textContent.Text), &responseData)
		require.NoError(t, err)

		// Either succeeded or failed with host not found
		assert.Contains(t, responseData, "failed")
	}
}

func TestIsHostNotFoundError(t *testing.T) {
	t.Run("DNSError with IsNotFound true", func(t *testing.T) {
		dnsErr := &net.DNSError{
			Err:         "no such host",
			Name:        "example.invalid",
			Server:      "8.8.8.8",
			IsNotFound:  true,
			IsTimeout:   false,
			IsTemporary: false,
		}

		assert.True(t, isHostNotFoundError(dnsErr))
	})

	t.Run("DNSError with IsNotFound false", func(t *testing.T) {
		dnsErr := &net.DNSError{
			Err:         "server misbehaving",
			Name:        "example.com",
			Server:      "8.8.8.8",
			IsNotFound:  false,
			IsTimeout:   false,
			IsTemporary: true,
		}

		assert.False(t, isHostNotFoundError(dnsErr))
	})

	t.Run("non-DNS error", func(t *testing.T) {
		err := context.DeadlineExceeded
		assert.False(t, isHostNotFoundError(err))
	})

	t.Run("nil error", func(t *testing.T) {
		assert.False(t, isHostNotFoundError(nil))
	})
}

func TestLookupIPAddresses(t *testing.T) {
	ctx := context.Background()

	t.Run("successful IPv4 lookup", func(t *testing.T) {
		addresses, err := lookupIPAddresses(ctx, "localhost", "ipv4")

		require.NoError(t, err)
		assert.NotEmpty(t, addresses)

		// Check that all addresses are valid IPv4
		for _, addr := range addresses {
			ip := net.ParseIP(addr)
			assert.NotNil(t, ip)
			assert.NotNil(t, ip.To4()) // Should be IPv4
		}
	})

	t.Run("successful IPv6 lookup", func(t *testing.T) {
		addresses, err := lookupIPAddresses(ctx, "localhost", "ipv6")
		if err != nil {
			// IPv6 might not be available on all systems
			t.Skipf("IPv6 lookup failed (might not be available): %v", err)
		}

		assert.NotEmpty(t, addresses)

		// Check that all addresses are valid IPv6
		for _, addr := range addresses {
			ip := net.ParseIP(addr)
			assert.NotNil(t, ip)
			assert.Nil(t, ip.To4()) // Should NOT be IPv4 (i.e., should be IPv6)
		}
	})

	t.Run("non-existent host", func(t *testing.T) {
		addresses, err := lookupIPAddresses(ctx, "this-does-not-exist-98765.invalid", "ipv4")

		require.Error(t, err)
		assert.Nil(t, addresses)
		assert.True(t, isHostNotFoundError(err))
	})
}
