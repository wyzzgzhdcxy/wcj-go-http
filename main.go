package main

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

func init() {
	if runtime.GOOS == "windows" {
		hideConsoleWindow()
	}
}

func hideConsoleWindow() {
	user32 := syscall.MustLoadDLL("user32.dll")
	ShowWindow := user32.MustFindProc("ShowWindow")
	kernel32 := syscall.MustLoadDLL("kernel32.dll")
	GetConsoleWindow := kernel32.MustFindProc("GetConsoleWindow")
	hwnd, _, _ := GetConsoleWindow.Call()
	if hwnd != 0 {
		ShowWindow.Call(hwnd, 0)
	}
}

func main() {
	port := ":38472"
	addr := "127.0.0.1" + port

	// 尝试关闭占用该端口的进程
	if err := killProcessOnPort(port); err != nil {
		fmt.Printf("端口清理: %v\n", err)
	}

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

		info, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				http.Error(w, "文件不存在: "+filePath, http.StatusNotFound)
			} else {
				http.Error(w, "无法访问: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		if info.IsDir() {
			http.Error(w, "暂不支持目录访问", http.StatusForbidden)
			return
		}

		fmt.Printf("访问文件: %s\n", filePath)
		w.Header().Set("Content-Type", detectContentType(filePath))
		http.ServeFile(w, r, filePath)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("文件服务器\n用法: /file?path=文件完整路径\n例: /file?path=E:/images/photo.jpg"))
	})

	fmt.Printf("服务器启动于 http://localhost%s\n", port)
	fmt.Println("仅允许本机访问...")
	http.ListenAndServe(addr, nil)
}

func killProcessOnPort(port string) error {
	ln, err := net.Listen("tcp", port)
	if err == nil {
		ln.Close()
		return nil
	}

	// 端口被占用，尝试获取 PID 并杀死
	proc1, _ := os.FindProcess(1)
	if proc1 != nil {
		// Windows: 使用 taskkill
		exec.Command("cmd", "/c", "for /f \"tokens=5\" %a in ('netstat -ano | findstr :"+strings.TrimPrefix(port, ":")+"') do taskkill /F /PID %a").Run()
	}

	// 等待端口释放
	for i := 0; i < 10; i++ {
		ln, err = net.Listen("tcp", port)
		if err == nil {
			ln.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return err
}

func detectContentType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".html", ".htm":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".webp":
		return "image/webp"
	case ".ico":
		return "image/x-icon"
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".avi":
		return "video/x-msvideo"
	case ".mkv":
		return "video/x-matroska"
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".ogg":
		return "audio/ogg"
	case ".flac":
		return "audio/flac"
	case ".pdf":
		return "application/pdf"
	case ".zip":
		return "application/zip"
	case ".tar":
		return "application/x-tar"
	case ".gz":
		return "application/gzip"
	case ".txt":
		return "text/plain"
	case ".xml":
		return "text/xml"
	case ".csv":
		return "text/csv"
	case ".doc", ".docx":
		return "application/msword"
	case ".xls", ".xlsx":
		return "application/vnd.ms-excel"
	default:
		return "application/octet-stream"
	}
}