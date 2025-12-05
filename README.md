# Reqstat

A HTTP request analyser CLI built with [Charm](https://charm.sh) libraries.

## Features

- 🚀 **Request timing breakdown** - DNS, TCP, TLS, server response
- 📊 **Visual timing bar** - See where time is spent
- 📦 **Response size analysis** - Content length and type  
- 📝 **Headers display** - All response headers, sorted
- 🔍 **JSON shape analysis** - Like `jq`, shows structure of JSON responses

## Installation

```bash
go install github.com/ewan-valentine/reqstat@latest
```

Or build from source:

```bash
git clone https://github.com/ewan-valentine/reqstat
cd reqstat
go build -o reqstat .
```

## Usage

```bash
# Basic GET request
reqstat get https://api.github.com/users/octocat

# Add custom headers
reqstat get https://api.example.com/data -H "Authorization: Bearer token"

# Show response body
reqstat get https://api.example.com/data --body

# Multiple headers
reqstat get https://api.example.com/data \
  -H "Authorization: Bearer token" \
  -H "Accept: application/json"
```

## Options

| Flag | Short | Description |
|------|-------|-------------|
| `--header` | `-H` | Add custom header (can be repeated) |
| `--body` | `-b` | Show response body |
| `--pretty` | `-p` | Pretty print JSON body (default: true) |
| `--max-body` | `-m` | Max body characters to display (default: 1000) |

## Example Output

```
⚡ reqstat
   GET https://api.github.com/users/octocat

│ STATUS

   ● 200 OK

│ TIMING

   Total              234ms
   DNS Lookup         12ms
   TCP Connect        45ms
   TLS Handshake      89ms
   Server Response    156ms

   ████████████████████████████████████████████████
   DNS │ TCP │ TLS │ Server │ Transfer

│ SIZE

   Content Length     1.2 KB
   Content Type       application/json; charset=utf-8

│ HEADERS

   Cache-Control: private, max-age=60
   Content-Type: application/json; charset=utf-8
   ...

│ JSON SHAPE

   Keys: 32 | Depth: 2 | Array items: 0

   {
     avatar_url: string // e.g. "https://avatars..."
     bio: string // e.g. "GitHub mascot"
     company: string // e.g. "@github"
     followers: number // e.g. 12345
     ...
   }
```

## License

MIT

