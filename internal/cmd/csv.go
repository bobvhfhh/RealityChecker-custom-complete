package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"

	"RealityChecker/internal/ui"
)

// executeCSV 从CSV文件批量检测域名
func (r *RootCmd) executeCSV(csvFile string) {
	// 检查文件是否存在
	if _, err := os.Stat(csvFile); os.IsNotExist(err) {
		ui.PrintErrorWithDetails(
			fmt.Sprintf("错误：CSV文件不存在 '%s'", csvFile),
			"请使用 RealiTLScanner 工具扫描，得到 CSV 文件",
			"命令：./RealiTLScanner -addr <VPS IP> -port 443 -thread 100 -timeout 5 -out file.csv",
			"（提示：RealiTLScanner 尽量在本地运行，不要在远端）",
			"（重要：多次运行时请更改输出文件名，如 file1.csv、file2.csv 等）",
		)
		return
	}

	// 读取CSV文件
	file, err := os.Open(csvFile)
	if err != nil {
		ui.PrintError(fmt.Sprintf("错误：无法打开CSV文件 '%s': %v", csvFile, err))
		return
	}
	defer file.Close()

	// 解析CSV
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		ui.PrintErrorWithDetails(
			fmt.Sprintf("错误：解析CSV文件失败: %v", err),
			"请使用 RealiTLScanner 工具扫描，得到 CSV 文件",
			"命令：./RealiTLScanner -addr <VPS IP> -port 443 -thread 100 -timeout 5 -out file.csv",
			"（提示：RealiTLScanner 尽量在本地运行，不要在远端）",
			"（重要：多次运行时请更改输出文件名，如 file1.csv、file2.csv 等）",
		)
		return
	}

	if len(records) < 2 {
		ui.PrintErrorWithDetails(
			"错误：CSV文件格式错误或为空",
			"请使用 RealiTLScanner 工具扫描，得到 CSV 文件",
			"命令：./RealiTLScanner -addr <VPS IP> -port 443 -thread 100 -timeout 5 -out file.csv",
			"（提示：RealiTLScanner 尽量在本地运行，不要在远端）",
			"（重要：多次运行时请更改输出文件名，如 file1.csv、file2.csv 等）",
		)
		return
	}

	// 提取域名（从CERT_DOMAIN列）
	domains := extractDomainsFromCSV(records)
	if len(domains) == 0 {
		ui.PrintErrorWithDetails(
			"错误：未找到有效的域名",
			"请使用 RealiTLScanner 工具扫描，得到 CSV 文件",
			"命令：./RealiTLScanner -addr <VPS IP> -port 443 -thread 100 -timeout 5 -out file.csv",
			"（提示：RealiTLScanner 尽量在本地运行，不要在远端）",
			"（重要：多次运行时请更改输出文件名，如 file1.csv、file2.csv 等）",
		)
		return
	}

	fmt.Printf("[%s] 从CSV文件提取到 %d 个域名\n", time.Now().Format("15:04:05"), len(domains))
	ui.PrintTimestampedMessage("开始批量检测...")

	_, err = r.batchManager.CheckDomains(r.ctx, domains)
	if err != nil {
		fmt.Printf("批量检测失败: %v\n", err)
		return
	}

	// 显示广告
	ui.PrintAdvertisement()
}

// extractDomainsFromCSV 从CSV记录中提取域名
func extractDomainsFromCSV(records [][]string) []string {
	var domains []string
	domainSet := make(map[string]bool)
	if len(records) < 2 {
		return domains
	}
	domainIndex, ok := findColumnIndex(records[0], "CERT_DOMAIN")
	if !ok {
		return domains
	}
	expectedFields := len(records[0])
	for i := 1; i < len(records); i++ {
		row, ok := normalizeCSVRow(records[i], domainIndex, expectedFields)
		if !ok {
			continue
		}
		certDomain := strings.TrimSpace(row[domainIndex])
		if certDomain == "" {
			continue
		}
		certDomain = strings.Trim(certDomain, string([]byte{34}))
		if shouldExcludeDomain(certDomain) {
			continue
		}
		if !domainSet[certDomain] {
			domains = append(domains, certDomain)
			domainSet[certDomain] = true
		}
	}
	return domains
}

func findColumnIndex(header []string, name string) (int, bool) {
	for i, column := range header {
		column = strings.TrimSpace(strings.TrimPrefix(column, "\ufeff"))
		if strings.EqualFold(column, name) {
			return i, true
		}
	}
	return -1, false
}

func normalizeCSVRow(row []string, domainIndex, expectedFields int) ([]string, bool) {
	if len(row) == expectedFields {
		return row, true
	}
	if len(row) < expectedFields || domainIndex >= len(row) {
		return nil, false
	}
	end := len(row) - (expectedFields - domainIndex - 1)
	if end <= domainIndex || end > len(row) {
		return nil, false
	}
	merged := append([]string{}, row[:domainIndex]...)
	merged = append(merged, strings.Join(row[domainIndex:end], ","))
	merged = append(merged, row[end:]...)
	return merged, len(merged) == expectedFields
}

// shouldExcludeDomain 判断是否应该排除某个域名
func shouldExcludeDomain(domain string) bool {
	// 1. 排除包含通配符(*)的域名
	if strings.Contains(domain, "*") {
		return true
	}

	// 2. 排除列表
	excludePatterns := []string{
		"localhost",
		"server.domain.com",
		"johnnasmalley.hostname",
		"Kubernetes Ingress Controller Fake Certificate",
		"CloudFlare Origin Certificate",
		"FortiGate",
		"Unspecified",
	}

	domainLower := strings.ToLower(domain)

	for _, pattern := range excludePatterns {
		if strings.Contains(domainLower, strings.ToLower(pattern)) {
			return true
		}
	}

	// 3. 排除IP地址格式
	if strings.Contains(domain, ".") && !strings.Contains(domain, "..") {
		parts := strings.Split(domain, ".")
		if len(parts) == 4 {
			// 可能是IP地址，简单检查
			isIP := true
			for _, part := range parts {
				if len(part) > 3 {
					isIP = false
					break
				}
			}
			if isIP {
				return true // 是IP地址，排除
			}
		}
	}

	// 4. 排除无效域名（太短或包含特殊字符）
	if len(domain) < 3 {
		return true
	}

	// 5. 排除包含多个连续点的域名
	if strings.Contains(domain, "..") {
		return true
	}

	return false
}
