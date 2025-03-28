# HTTP 方法覆盖检测工具
[English](README.md) | 简体中文

## 📖 项目简介
一款专注于 HTTP 方法覆盖检测的安全工具，采用渐进式探测策略，最大限度降低对目标系统的影响。基于 [OWASP WSTG-CONF-06](https://owasp.org/www-project-web-security-testing-guide/latest/4-Web_Application_Security_Testing/02-Configuration_and_Deployment_Management_Testing/06-Test_HTTP_Methods) 测试指南设计，用于识别和验证 HTTP 方法覆盖漏洞。

## ⚙️ 工作原理
该工具采用三阶段检测策略：

### 1️⃣ OPTIONS 可用性检测
- 首先尝试 OPTIONS 方法获取服务器支持的 HTTP 方法列表
- 如果 OPTIONS 方法不可用，转入方法覆盖测试

### 2️⃣ 方法覆盖测试
- 当 OPTIONS 不可用时，尝试通过方法覆盖的方式获取 OPTIONS 信息
- 使用多种标准的 HTTP 方法覆盖请求头进行测试

### 3️⃣ 安全性验证
- 基于服务器返回的允许方法列表
- 优先选择安全性高的方法进行覆盖测试
- 避免使用具有破坏性的方法（如 DELETE）

## 🚀 快速开始

### 安装
```bash
git clone https://github.com/kingjly/http-override.git
cd http-override
go build
```

### 使用示例
```bash
# 扫描单个目标
./http_override -u https://example.com

# 批量扫描
./http_override -l urls.txt -c 5 -t 10
```

## 📝 命令行参数
| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-u` | 指定单个目标 URL | - |
| `-l` | 指定 URL 列表文件 | - |
| `-c` | 并发数 | 5 |
| `-t` | 超时时间(秒) | 10 |

## 🛠️ 支持的方法覆盖请求头
- `X-HTTP-Method-Override`
- `X-HTTP-Method`
- `X-Method-Override`
- `_method`
- `X-Original-HTTP-Method`
- `X-Override-Method`

## ⚠️ 注意事项
1. 本工具采用渐进式探测策略，优先使用影响较小的检测方法
2. 仅用于授权的安全测试，请勿用于未授权的测试活动
3. 建议在测试环境中先进行验证

## 📄 许可证
[MIT License](LICENSE)
