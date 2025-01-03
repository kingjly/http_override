package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"http_override/scanner"
)

func main() {
	var (
		url     string
		urlFile string
		workers int
		timeout int
	)

	flag.StringVar(&url, "u", "", "目标URL")
	flag.StringVar(&urlFile, "l", "", "包含URL列表的文件")
	flag.IntVar(&workers, "c", 5, "并发数 (默认: 5)")
	flag.IntVar(&timeout, "t", 10, "超时时间(秒) (默认: 10)")
	flag.Parse()

	if url == "" && urlFile == "" {
		fmt.Println("请使用 -u 指定单个URL 或 -l 指定URL列表文件")
		fmt.Println("\n用法:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if workers < 1 {
		fmt.Println("并发数必须大于0")
		os.Exit(1)
	}

	httpScanner := scanner.New(time.Duration(timeout) * time.Second)

	if url != "" {
		httpScanner.Scan(url)
		return
	}

	file, err := os.Open(urlFile)
	if err != nil {
		fmt.Printf("打开文件错误: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var wg sync.WaitGroup
	sem := make(chan struct{}, workers)

	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		url := strings.TrimSpace(fileScanner.Text())
		if url == "" {
			continue
		}

		wg.Add(1)
		sem <- struct{}{}

		go func(url string) {
			defer wg.Done()
			defer func() { <-sem }()
			httpScanner.Scan(url)
		}(url)
	}

	wg.Wait()

	if err := fileScanner.Err(); err != nil {
		fmt.Printf("读取文件错误: %v\n", err)
		os.Exit(1)
	}
}
