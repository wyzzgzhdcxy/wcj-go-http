# Static File Server

本地静态文件服务器，通过 HTTP 协议访问电脑上的任意文件。

## 功能特性

- 访问电脑任意路径的文件（图片、视频、音频、文档等）
- 自动识别文件类型，设置正确的 Content-Type
- 只监听本机（127.0.0.1），外网无法访问
- 再次启动自动关闭旧进程
- Windows 下无黑窗口
- 支持中文路径、特殊符号、Emoji

## 使用方法

### 启动服务器

```bash
./static-server.exe
```

服务启动后访问地址：`http://localhost:38472`

### 访问文件

```
http://localhost:38472/file?path=文件完整路径
```

**示例**：
```
http://localhost:38472/file?path=E:/images/photo.jpg
http://localhost:38472/file?path=E:/video/movie.mp4
http://localhost:38472/file?path=C:/Users/documents/report.pdf
```

### 路径格式

- Windows 路径使用正斜杠：`E:/images/photo.jpg`
- 特殊字符（+、空格、& 等）会自动处理
- 中文路径直接写即可

## 文件类型支持

| 类型 | 扩展名 |
|------|--------|
| 图片 | png, jpg, gif, svg, webp, ico |
| 视频 | mp4, webm, avi, mkv |
| 音频 | mp3, wav, ogg, flac |
| 文档 | pdf, doc, xls |
| 压缩 | zip, tar, gz |
| 其他 | html, css, js, json, txt, xml, csv |

## 编译打包

### Windows (最小 GUI 包)

```bash
go build -ldflags="-H windowsgui -s -w" -o static-server.exe main.go
```

### 验证

```bash
# 启动服务
./static-server.exe

# 测试访问
curl http://localhost:38472/file?path=main.go
```

## 项目结构

```
.
├── main.go           # 主程序（包含 Windows 隐藏窗口逻辑）
└── static-server.exe # 编译后的可执行文件
```

## 注意事项

1. 路径中的 `+` 符号在 URL 中会表示空格，如需使用可使用 `%2B` 代替
2. 服务器只允许本机访问，请勿在公网环境使用
3. 文件夹目录暂不支持访问
