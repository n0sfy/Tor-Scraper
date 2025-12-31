package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

func main() {

	proxyAddress := "127.0.0.1:9050"

	dialer, err := proxy.SOCKS5("tcp", proxyAddress, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	transport := &http.Transport{
		Dial:              dialer.Dial,
		DisableKeepAlives: true,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	logFile, err := os.OpenFile(
		"scan_report.log",
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)

	fmt.Println("[INFO] Checking tor connection...")
	resp, err := client.Get("https://check.torproject.org/api/ip")
	if err != nil {
		log.Fatalf("[ERR] Connection failed: %v", err)
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("[INFO] Current IP Status: %s\n", string(body))
		resp.Body.Close()
	}

	file, err := os.Open("targets.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	outputDir := "Results"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, 0755)
	}

	scanner := bufio.NewScanner(file)
	fmt.Println("[INFO] Scanning started...")

	for scanner.Scan() {
		satir := scanner.Text()
		url := strings.TrimSpace(satir)

		if url == "" {
			continue
		}

		fmt.Printf("[INFO] Scanning: %s ... ", url)

		resp, err := client.Get(url)
		if err != nil {
			fmt.Printf("[ERR] TIMEOUT/FAIL (%v)\n", err)
			logger.Printf("FAIL %s ERROR=%v", url, err)

			continue
		}

		if resp.StatusCode != 200 {
			fmt.Printf("[FAIL] HTTP Status: %s\n", resp.Status)
			logger.Printf("FAIL %s STATUS=%s", url, resp.Status)
			resp.Body.Close()
			continue
		}

		htmlData, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			fmt.Printf("[ERR] HTML could not be read\n")
		} else {
			safeName := strings.ReplaceAll(url, "://", "_")
			safeName = strings.ReplaceAll(safeName, "/", "_")
			safeName = strings.ReplaceAll(safeName, ":", "_")

			timestamp := time.Now().Unix()

			filename := fmt.Sprintf("%s/%s_%d.html", outputDir, safeName, timestamp)

			writeErr := os.WriteFile(filename, htmlData, 0644)
			if writeErr != nil {
				fmt.Printf("[ERR] Registration failed: %v\n", writeErr)
			} else {
				fmt.Printf("[SUCCESS] Saved: %s\n", filename)
				logger.Printf("SUCCESS %s FILE=%s STATUS=%s", url, filename, resp.Status)

			}
		}

		resp.Body.Close()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
