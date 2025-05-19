# DNS and WHOIS query MCP server `mcp-domaintools`

`mcp-domaintools` is a [Model Context Protocol (MCP)](https://modelcontextprotocol.io/introduction) server providing DNS and WHOIS query capabilities for AI assistants. It enables AI models to perform DNS lookups both via local DNS resolvers and remote DNS-over-HTTPS services.

For local DNS queries, it uses the system's configured DNS servers. For remote DNS queries, it uses Cloudflare DNS-over-HTTPS queries with a fallback to Google DNS-over-HTTPS. This is more than enough for most use cases.

For custom DNS-over-HTTPS servers, you can use the `--remote-server-address` flag. The server endpoint must implement the HTTP reponse format as defined by [RFC 8484](https://datatracker.ietf.org/doc/html/rfc8484#section-4.2).

For custom WHOIS servers, you can use the `--custom-whois-server` flag. The server endpoint must implement the HTTP reponse format as defined by [RFC 3912](https://datatracker.ietf.org/doc/html/rfc3912), although plain text responses are also supported.

## Features

- **Local DNS Queries**: Perform DNS lookups using the OS-configured DNS servers
- **Remote DNS-over-HTTPS**: Perform secure DNS queries via Cloudflare and Google DNS-over-HTTPS services
- **WHOIS Lookups**: Perform WHOIS queries to get domain registration information
- **Multiple Record Types**: Support for A, AAAA, CNAME, MX, NS, PTR, SOA, SRV, and TXT record types
- **Fallback Mechanism**: Automatically tries multiple DNS servers for reliable results

## Installation

### Editor Configuration

Add the following configuration to your editor's settings to use `mcp-domaintools`:

```json5
{
  "mcpServers": {
    "dns": {
      "command": "mcp-domaintools",
      "args": [
        // Uncomment and modify as needed:
        // "--remote-server-address=https://your-custom-doh-server.com/dns-query",
        // "--custom-whois-server=whois.yourdomain.com",
        // "--timeout=10s"
      ],
      "env": {}
    }
  }
}
```

You can use `mcp-domaintools` directly from your `$PATH` as shown above, or provide the full path to the binary (e.g., `/path/to/mcp-domaintools`).

Alternatively, you can run `mcp-domaintools` directly with Docker without installing the binary:

```json5
{
  "mcpServers": {
    "dns": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "ghcr.io/patrickdappollonio/mcp-domaintools:latest"
        // Add custom options if needed:
        // "--remote-server-address=https://your-custom-doh-server.com/dns-query",
        // "--custom-whois-server=whois.yourdomain.com",
        // "--timeout=10s"
      ],
      "env": {}
    }
  }
}
```

See ["Available MCP Tools"](#available-mcp-tools) for information on the tools exposed by `mcp-domaintools`.

### Homebrew (macOS and Linux)

```bash
brew install patrickdappollonio/tap/mcp-domaintools
```

### Docker

The MCP server is available as a Docker image using `stdio` to communicate:

```bash
docker pull ghcr.io/patrickdappollonio/mcp-domaintools:latest
docker run --rm ghcr.io/patrickdappollonio/mcp-domaintools:latest
```

Check the implementation above on how to configure the MCP server to run as a container in your editor or tool.

### GitHub Releases

Download the pre-built binaries for your platform from the [GitHub Releases page](https://github.com/patrickdappollonio/mcp-domaintools/releases).

## Available MCP Tools

There are 3 tools available:

- `local_dns_query`: Perform DNS queries against the local DNS resolver as configured by the OS
- `remote_dns_query`: Perform DNS queries against a remote DNS-over-HTTPS server
- `whois_query`: Perform WHOIS lookups to get domain registration information

### Local DNS Query

Performs DNS queries using local OS-defined DNS servers.

**Arguments:**
- `domain` (required): The domain name to query (e.g., example.com)
- `record_type` (required): Type of DNS record to query (A, AAAA, CNAME, MX, NS, PTR, SOA, SRV, TXT)

### Remote DNS Query

Performs DNS queries using remote DNS-over-HTTPS servers (Google and Cloudflare).

**Arguments:**
- `domain` (required): The domain name to query (e.g., example.com)
- `record_type` (required): Type of DNS record to query (A, AAAA, CNAME, MX, NS, PTR, SOA, SRV, TXT)

### WHOIS Query

Performs WHOIS lookups to get domain registration information.

**Arguments:**
- `domain` (required): The domain name to query (e.g., example.com)
