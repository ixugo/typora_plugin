package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("用法: go run main.go <HTML文件>")
		os.Exit(1)
	}

	filename := os.Args[1]
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("读取文件错误: %v\n", err)
		os.Exit(1)
	}

	htmlContent := string(content)
	htmlDir := filepath.Dir(filename)
	re := regexp.MustCompile(`<img[^>]*src="([^"]+)"[^>]*>`)

	var i int
	htmlContent = re.ReplaceAllStringFunc(htmlContent, func(match string) string {
		src := re.FindStringSubmatch(match)[1]
		if strings.HasPrefix(src, "data:") {
			return match
		}

		i++
		mimeType := "image/png"
		switch filepath.Ext(src) {
		case ".jpg", ".jpeg":
			mimeType = "image/jpeg"
		case ".gif":
			mimeType = "image/gif"
		case ".svg":
			mimeType = "image/svg+xml"
		}

		var base64Img string
		var err error
		if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
			base64Img, err = fetchAndEncodeImage(src)
		} else {
			base64Img, err = localFileToBase64(filepath.Join(htmlDir, src))
		}

		if err != nil {
			fmt.Printf("图片转换失败: %s, 错误: %v\n", src, err)
			return match
		}

		newSrc := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Img)
		return strings.Replace(match, src, newSrc, 1)
	})

	if err := os.WriteFile(filename, []byte(htmlContent), 0o600); err != nil {
		fmt.Printf("写入文件错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("共 %d 处修改，图片转换完成: %s\n", i, filename)
}

func fetchAndEncodeImage(link string) (string, error) {
	u, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	resp, err := http.Get(u.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP请求失败，状态码: %d", resp.StatusCode)
	}

	imgBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(imgBytes), nil
}

func localFileToBase64(path string) (string, error) {
	imgBytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(imgBytes), nil
}
