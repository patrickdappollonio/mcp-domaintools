# Network and domain tools MCP server `mcp-domaintools`

<img src="https://i.imgur.com/cai3zrG.png" width="160" align="right" />  `mcp-domaintools` is a [Model Context Protocol (MCP)](https://modelcontextprotocol.io/introduction) server providing comprehensive network and domain analysis capabilities for AI assistants. It enables AI models to perform DNS lookups, WHOIS queries, connectivity testing, TLS certificate analysis, and hostname resolution.

For local DNS queries, it uses the system's configured DNS servers. For remote DNS queries, it uses Cloudflare DNS-over-HTTPS queries with a fallback to Google DNS-over-HTTPS. This is more than enough for most use cases.

For custom DNS-over-HTTPS servers, you can use the `--remote-server-address` flag. The server endpoint must implement the HTTP reponse format as defined by [RFC 8484](https://datatracker.ietf.org/doc/html/rfc8484#section-4.2).

For custom WHOIS servers, you can use the `--custom-whois-server` flag. The server endpoint must implement the HTTP reponse format as defined by [RFC 3912](https://datatracker.ietf.org/doc/html/rfc3912), although plain text responses are also supported.

## Features

- **Local DNS Queries**: Perform DNS lookups using the OS-configured DNS servers
- **Remote DNS-over-HTTPS**: Perform secure DNS queries via Cloudflare and Google DNS-over-HTTPS services
- **WHOIS Lookups**: Perform WHOIS queries to get domain registration information
- **Hostname Resolution**: Convert hostnames to their corresponding IP addresses (IPv4, IPv6, or both)
- **Ping Operations**: Test connectivity and measure response times to hosts
- **TLS Certificate Analysis**: Check TLS certificate chains for validity, expiration, and detailed certificate information
- **Multiple Record Types**: Support for A, AAAA, CNAME, MX, NS, PTR, SOA, SRV, and TXT record types
- **Fallback Mechanism**: Automatically tries multiple DNS servers for reliable results
- **SSE Support**: Run as an HTTP server with Server-Sent Events (SSE) for web-based integrations

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
        // "--timeout=10s",
        // "--ping-timeout=5s",
        // "--ping-count=4",
        // "--tls-timeout=10s"
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
        // "--timeout=10s",
        // "--ping-timeout=5s",
        // "--ping-count=4",
        // "--tls-timeout=10s"
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

For SSE mode with Docker, expose the SSE port (default `3000`):

```bash
docker run --rm -p 3000:3000 ghcr.io/patrickdappollonio/mcp-domaintools:latest --sse --sse-port 3000
```

Check the implementation above on how to configure the MCP server to run as a container in your editor or tool.

### Cursor

You can use one-click to install in Cursor (note this will use the Docker version of the MCP server since it doesn't require a local binary installation):

[![Install MCP Server](assets/cursor.svg)](cursor://anysphere.cursor-deeplink/mcp/install?name=domaintools&config=eyJjb21tYW5kIjoiZG9ja2VyIiwiYXJncyI6WyJydW4iLCItaSIsIi0tcm0iLCJnaGNyLmlvL3BhdHJpY2tkYXBwb2xsb25pby9tY3AtZG9tYWludG9vbHM6bGF0ZXN0Il0sImVudiI6e319)

### GitHub Releases

Download the pre-built binaries for your platform from the [GitHub Releases page](https://github.com/patrickdappollonio/mcp-domaintools/releases).

## Available MCP Tools

There are 6 tools available:

- `local_dns_query`: Perform DNS queries against the local DNS resolver as configured by the OS
- `remote_dns_query`: Perform DNS queries against a remote DNS-over-HTTPS server
- `whois_query`: Perform WHOIS lookups to get domain registration information
- `resolve_hostname`: Convert a hostname to its corresponding IP addresses (IPv4, IPv6, or both)
- `ping`: Perform ping operations to test connectivity and measure response times to a host
- `tls_certificate_check`: Check TLS certificate chain for a domain to analyze certificate validity, expiration, and chain structure

## Running Modes

### Standard (stdio) Mode

By default, `mcp-domaintools` runs in stdio mode, which is suitable for integration with editors and other tools that communicate via standard input/output.

```bash
mcp-domaintools
```

### Server-Sent Events (SSE) Mode

Alternatively, you can run `mcp-domaintools` as an HTTP server with SSE support for web-based integrations:

```bash
mcp-domaintools --sse --sse-port=3000
```

In SSE mode, the server will listen on the specified port (default: 3000) and provide the same MCP tools over HTTP using Server-Sent Events. This is useful for web applications or environments where stdio communication isn't practical.

## Configuration Options

The following command-line flags are available to configure the MCP server:

**General Options:**
- `--timeout=DURATION`: Timeout for DNS queries (default: 5s)
- `--remote-server-address=URL`: Custom DNS-over-HTTPS server address
- `--custom-whois-server=ADDRESS`: Custom WHOIS server address

**Ping Options:**
- `--ping-timeout=DURATION`: Timeout for ping operations (default: 5s)
- `--ping-count=NUMBER`: Default number of ping packets to send (default: 4)

**TLS Options:**
- `--tls-timeout=DURATION`: Timeout for TLS certificate checks (default: 10s)

**SSE Server Options:**
- `--sse`: Enable SSE server mode
- `--sse-port=PORT`: Specify the port to listen on (default: 3000)

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

### Hostname Resolution

Converts a hostname to its corresponding IP addresses.

**Arguments:**
- `hostname` (required): The hostname to resolve (e.g., example.com)
- `ip_version` (optional): IP version to resolve (ipv4, ipv6, or both); defaults to ipv4

### Ping

Performs ping operations to test connectivity and measure response times to a host.

**Arguments:**
- `target` (required): The hostname or IP address to ping (e.g., example.com or 8.8.8.8)
- `count` (optional): Number of ping packets to send; defaults to 4

### TLS Certificate Check

Checks TLS certificate chain for a domain to analyze certificate validity, expiration, and chain structure.

**Arguments:**
- `domain` (required): The domain name to check TLS certificate for (e.g., example.com)
- `port` (optional): Port to connect to for TLS check; defaults to 443
- `include_chain` (optional): Whether to include the full certificate chain in the response; defaults to true
- `check_expiry` (optional): Whether to check certificate expiration and provide warnings; defaults to true
- `server_name` (optional): Server name for SNI (Server Name Indication); defaults to the domain name
