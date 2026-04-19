package main

import (
	"fmt"
	"mime"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	port := ":38472"
	addr := "127.0.0.1" + port
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("端口 %s 被占用: %v\n", port, err)
		os.Exit(1)
	}
	ln.Close()

	http.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		var filePath string

		if r.URL.RawQuery != "" {
			params, _ := url.ParseQuery(r.URL.RawQuery)
			filePath = params.Get("path")
		}

		if filePath == "" {
			http.Error(w, "请提供 path 参数\n例: /file?path=E:/images/photo.jpg", http.StatusBadRequest)
			return
		}

		filePath = strings.ReplaceAll(filePath, "\\", "/")
		cleanPath := filepath.Clean(filePath)

		if strings.Contains(cleanPath, "..") {
			http.Error(w, "非法路径", http.StatusForbidden)
			return
		}

		info, err := os.Stat(cleanPath)
		if err != nil {
			if os.IsNotExist(err) {
				http.Error(w, "文件不存在: "+cleanPath, http.StatusNotFound)
			} else {
				http.Error(w, "无法访问: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		if info.IsDir() {
			http.Error(w, "暂不支持目录访问", http.StatusForbidden)
			return
		}

		ext := strings.ToLower(filepath.Ext(cleanPath))
		contentType := mime.TypeByExtension(ext)
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		w.Header().Set("Content-Type", contentType)
		http.ServeFile(w, r, cleanPath)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("文件服务器\n用法: /file?path=文件完整路径\n例: /file?path=E:/images/photo.jpg"))
	})

	fmt.Printf("文件服务器启动于 http://localhost%s\n", port)
	http.ListenAndServe(addr, nil)
}
