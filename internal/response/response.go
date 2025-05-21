package response

import (
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// JSON encodes the provided data to JSON and wraps it in an MCP tool result.
// It returns the tool result and any error encountered during JSON marshaling.
func JSON(data interface{}) (*mcp.CallToolResult, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error generating JSON: %w", err)
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}
