# HTTP Method Override Detection Tool
简体中文 | [English](README.md)

## 📖 Introduction
A security tool focused on HTTP method override detection, using a progressive detection strategy to minimize impact on target systems. Designed based on [OWASP WSTG-CONF-06](https://owasp.org/www-project-web-security-testing-guide/latest/4-Web_Application_Security_Testing/02-Configuration_and_Deployment_Management_Testing/06-Test_HTTP_Methods) testing guidelines for identifying and verifying HTTP method override vulnerabilities.

## ⚙️ Working Principle
The tool employs a three-phase detection strategy:

### 1️⃣ OPTIONS Availability Detection
- First attempts to use the OPTIONS method to obtain the server's supported HTTP methods
- If OPTIONS is unavailable, proceeds to method override testing

### 2️⃣ Method Override Testing
- When OPTIONS is unavailable, attempts to obtain OPTIONS information through method override
- Tests using multiple standard HTTP method override headers

### 3️⃣ Security Verification
- Based on the server's allowed methods list
- Prioritizes safer methods for override testing
- Avoids destructive methods (such as DELETE)

## 🚀 Quick Start

### Installation
```bash
git clone https://github.com/yourusername/http-override.git
cd http-override
go build
```

### Usage Examples
```bash
# Scan single target
./http_override -u https://example.com

# Batch scanning
./http_override -l urls.txt -c 5 -t 10
```

## 📝 Command Line Parameters
| Parameter | Description | Default |
|-----------|-------------|---------|
| `-u` | Specify single target URL | - |
| `-l` | Specify URL list file | - |
| `-c` | Concurrency | 5 |
| `-t` | Timeout (seconds) | 10 |

## 🛠️ Supported Method Override Headers
- `X-HTTP-Method-Override`
- `X-HTTP-Method`
- `X-Method-Override`
- `_method`
- `X-Original-HTTP-Method`
- `X-Override-Method`

## ⚠️ Precautions
1. This tool uses a progressive detection strategy, prioritizing less impactful detection methods
2. For authorized security testing only, do not use for unauthorized testing activities
3. Recommended to verify in a test environment first

## 📄 License
[MIT License](LICENSE)
