# HTTP 方法覆盖漏洞扫描器

中文|[English](README.md) 

## 项目简介
HTTP 方法覆盖漏洞扫描器是一个用 Go 语言开发的安全工具，专门用于检测 Web 服务器中的 HTTP 方法覆盖（Method Override）漏洞。

## 功能特点
- 🔍 检测 OPTIONS 方法覆盖漏洞
- 🕵️ 检测 TRACE 方法覆盖漏洞
- 🚀 支持单 URL 和批量 URL 扫描
- ⚙️ 可配置并发数和超时时间
- 🎨 彩色控制台输出，便于阅读

## 安装方法
```bash
git clone https://github.com/yourusername/http_override.git
cd http_override
go build
```

## 使用示例
### 扫描单个 URL
```bash
./http_override -u https://example.com
```

### 批量扫描 URL
```bash
./http_override -l urls.txt
```

## 参数说明
| 参数 | 描述 | 默认值 |
|------|------|--------|
| `-u` | 指定单个目标 URL | 无 |
| `-l` | 指定包含 URL 列表的文件 | 无 |
| `-c` | 设置并发数 | 5 |
| `-t` | 设置超时时间（秒） | 10 |

## 注意事项
⚠️ 仅用于安全测试和研究，请确保获得授权后再对目标进行扫描

## 许可证
[MIT License](LICENSE)

## 贡献
欢迎提交 Issue 和 Pull Request！
