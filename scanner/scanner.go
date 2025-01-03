package scanner

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var methodOverrideHeaders = []string{
	"X-HTTP-Method-Override",
	"X-HTTP-Method",
	"X-Method-Override",
	"_method",
	"X-Original-HTTP-Method",
	"X-Override-Method",
}

type Result struct {
	StatusCode     int
	Headers        map[string][]string
	ResponseBody   string
	ContentType    string
	ContentLength  int64
	OverrideHeader string
}

type Scanner struct {
	client  *http.Client
	timeout time.Duration
}

func New(timeout time.Duration) *Scanner {
	tr := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}
	return &Scanner{
		client:  client,
		timeout: timeout,
	}
}

func (s *Scanner) checkOptions(url string) (*Result, error) {
	req, err := http.NewRequest(http.MethodOptions, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "*/*")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := &Result{
		StatusCode: resp.StatusCode,
		Headers:    make(map[string][]string),
	}

	for _, header := range []string{"Allow", "Access-Control-Allow-Methods"} {
		if values := resp.Header[header]; len(values) > 0 {
			result.Headers[header] = values
		}
	}

	return result, nil
}

func (s *Scanner) testMethodOverride(url string, baseMethod string, targetMethod string) ([]*Result, error) {
	var results []*Result

	for _, overrideHeader := range methodOverrideHeaders {
		req, err := http.NewRequest(baseMethod, url, nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", "Mozilla/5.0")
		req.Header.Set("Accept", "*/*")
		req.Header.Set(overrideHeader, targetMethod)

		if targetMethod == "TRACE" {
			req.Header.Set("X-Test-Trace", "test-value")
		}

		resp, err := s.client.Do(req)
		if err != nil {
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}

		result := &Result{
			StatusCode:     resp.StatusCode,
			Headers:        make(map[string][]string),
			ResponseBody:   string(body),
			ContentType:    resp.Header.Get("Content-Type"),
			ContentLength:  resp.ContentLength,
			OverrideHeader: overrideHeader,
		}

		for k, v := range resp.Header {
			result.Headers[k] = v
		}

		results = append(results, result)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("所有请求头测试都失败")
	}

	return results, nil
}

func isValidTraceResponse(result *Result) bool {
	if result.StatusCode != 200 {
		return false
	}

	contentType := strings.ToLower(result.ContentType)
	if !strings.Contains(contentType, "message/http") {
		return false
	}

	responseBody := strings.ToLower(result.ResponseBody)
	expectedHeaders := []string{
		strings.ToLower(result.OverrideHeader),
		"user-agent",
		"accept",
		"x-test-trace",
	}

	matchCount := 0
	for _, header := range expectedHeaders {
		if strings.Contains(responseBody, header) {
			matchCount++
		}
	}

	return matchCount >= 3
}

func getAllowedMethods(headers map[string][]string) []string {
	var methods []string

	if allow, ok := headers["Allow"]; ok && len(allow) > 0 {
		methods = append(methods, strings.Split(allow[0], ",")...)
	}

	if cors, ok := headers["Access-Control-Allow-Methods"]; ok && len(cors) > 0 {
		methods = append(methods, strings.Split(cors[0], ",")...)
	}

	for i, method := range methods {
		methods[i] = strings.TrimSpace(method)
	}

	return unique(methods)
}

func unique(slice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(strings.TrimSpace(s), item) {
			return true
		}
	}
	return false
}

func (s *Scanner) testTraceOverride(url, safeBase string) {
	PrintInfo("正在测试 HTTP 方法覆盖漏洞")
	fmt.Printf("    基础请求方法: %s\n", safeBase)
	fmt.Printf("    目标覆盖方法: TRACE\n")

	results, err := s.testMethodOverride(url, safeBase, "TRACE")
	if err != nil {
		PrintError("测试失败: %v", err)
		return
	}

	foundVulnerable := false
	for _, result := range results {
		if isValidTraceResponse(result) {
			foundVulnerable = true
			PrintVulnFound("发现 HTTP 方法覆盖漏洞!")
			fmt.Printf("\n%s%s漏洞详情:%s\n", Bold, ColorYellow, ColorReset)
			PrintVulnDetail("URL", url)
			PrintVulnDetail("基础方法", safeBase)
			PrintVulnDetail("目标方法", "TRACE")
			PrintVulnDetail("覆盖请求头", result.OverrideHeader)
			PrintVulnDetail("响应状态", fmt.Sprintf("%d", result.StatusCode))
			PrintVulnDetail("Content-Type", result.ContentType)
			PrintVulnDetail("响应长度", fmt.Sprintf("%d", result.ContentLength))

			if len(result.ResponseBody) > 0 {
				fmt.Printf("\n%s%s响应内容预览:%s\n", Bold, ColorYellow, ColorReset)
				preview := result.ResponseBody
				if len(preview) > 200 {
					preview = preview[:200] + "..."
				}
				fmt.Printf("%s%s%s\n", ColorCyan, preview, ColorReset)
			}
			fmt.Println()
		}
	}

	if !foundVulnerable {
		PrintSuccess("未发现方法覆盖漏洞")
		fmt.Printf("    测试结果摘要:\n")
		for _, result := range results {
			fmt.Printf("    - %s: %d\n", result.OverrideHeader, result.StatusCode)
		}
	}
}

func (s *Scanner) testOptionsOverride(url string) {
	results, err := s.testMethodOverride(url, "GET", "OPTIONS")
	if err != nil {
		PrintError("测试失败: %v", err)
		return
	}

	foundVulnerable := false
	for _, result := range results {
		if result.StatusCode == http.StatusOK {
			foundVulnerable = true
			PrintVulnFound("发现方法覆盖漏洞! (OPTIONS via GET)")
			fmt.Printf("\n%s%s漏洞详情:%s\n", Bold, ColorYellow, ColorReset)
			PrintVulnDetail("URL", url)
			PrintVulnDetail("覆盖请求头", result.OverrideHeader)
			PrintVulnDetail("响应状态", fmt.Sprintf("%d", result.StatusCode))
			fmt.Println()
		}
	}

	if !foundVulnerable {
		PrintSuccess("未发现方法覆盖漏洞")
		fmt.Printf("    测试结果摘要:\n")
		for _, result := range results {
			fmt.Printf("    - %s: %d\n", result.OverrideHeader, result.StatusCode)
		}
	}
}

func (s *Scanner) Scan(url string) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	PrintInfo("测试目标: %s", url)

	optionsResult, err := s.checkOptions(url)
	if err != nil {
		PrintError("OPTIONS 请求失败: %v", err)
		return
	}

	if optionsResult.StatusCode != http.StatusOK {
		PrintInfo("OPTIONS 方法返回状态码: %d", optionsResult.StatusCode)
		s.testOptionsOverride(url)
		return
	}

	PrintSuccess("OPTIONS 方法可用 (状态码: %d)", optionsResult.StatusCode)
	allowedMethods := getAllowedMethods(optionsResult.Headers)
	if len(allowedMethods) == 0 {
		PrintError("未能获取到允许的方法列表")
		return
	}

	PrintSuccess("服务器允许的方法: %s", strings.Join(allowedMethods, ", "))

	safeBase := "GET"
	if contains(allowedMethods, "POST") {
		safeBase = "POST"
	}

	if !contains(allowedMethods, "TRACE") {
		s.testTraceOverride(url, safeBase)
	}
}
