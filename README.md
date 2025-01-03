# HTTP Method Override Vulnerability Scanner

中文 | [English](README.md)

## Project Overview
HTTP Method Override Vulnerability Scanner is a security tool developed in Go, specifically designed to detect HTTP method override vulnerabilities in web servers.

## Features
- 🔍 Detect OPTIONS method override vulnerabilities
- 🕵️ Detect TRACE method override vulnerabilities
- 🚀 Support single URL and batch URL scanning
- ⚙️ Configurable concurrency and timeout
- 🎨 Colorful console output for easy reading

## Installation
```bash
git clone https://github.com/yourusername/http_override.git
cd http_override
go build
```

## Usage Examples
### Scan a Single URL
```bash
./http_override -u https://example.com
```

### Batch URL Scanning
```bash
./http_override -l urls.txt
```

## Parameter Description
| Parameter | Description | Default |
|----------|-------------|---------|
| `-u` | Specify a single target URL | None |
| `-l` | Specify a file containing URL list | None |
| `-c` | Set concurrency | 5 |
| `-t` | Set timeout in seconds | 10 |

## Precautions
⚠️ For security testing and research only. Ensure authorization before scanning targets.

## License
[MIT License](LICENSE)

## Contribution
Issues and Pull Requests are welcome!
