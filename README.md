# RealityChecker

RealityChecker 是一个用于评估 Reality 协议目标域名的命令行工具。它支持单域名、批量域名和 RealiTLScanner CSV 批量检测。

## 主要改进

本版本兼容 RealiTLScanner 的原始 CSV 输出：

- 按 `CERT_DOMAIN` 表头定位域名列，不依赖固定列号；
- 自动去重域名；
- 兼容域名字段中未正确引用的逗号，并显示修复/跳过的异常行数量；
- 支持带 UTF-8 BOM 的 CSV 表头；
- 缺少 `CERT_DOMAIN` 时给出明确错误。

## Linux VPS 部署

### 方式一：使用 Release（推荐）

从仓库的 [Releases](https://github.com/bobvhfhh/RealityChecker-custom-complete/releases) 下载对应架构的压缩包：

```bash
# x86_64 / amd64
wget -O reality-checker.zip https://github.com/bobvhfhh/RealityChecker-custom-complete/releases/latest/download/reality-checker-linux-amd64.zip
unzip reality-checker.zip
chmod +x reality-checker
./reality-checker version
```

ARM64 VPS 请下载 `reality-checker-linux-arm64.zip`。

### 方式二：从源码编译

需要 Go 1.21 或更高版本：

```bash
git clone https://github.com/bobvhfhh/RealityChecker-custom-complete.git
cd RealityChecker-custom-complete
go build -o reality-checker .
./reality-checker version
```

程序不是 Web 服务，不会监听端口，也不需要 Docker、systemd 或反向代理。

## 使用方法

### 单域名

```bash
./reality-checker check apple.com
```

### 多域名

```bash
./reality-checker batch apple.com tesla.com microsoft.com
```

### RealiTLScanner CSV

先用 RealiTLScanner 生成 CSV：

```bash
./RealiTLScanner -addr <VPS IP> -port 443 -thread 100 -timeout 5 -out file1.csv
```

然后直接检测原始 CSV：

```bash
./reality-checker csv file1.csv
./reality-checker csv file2.csv
```

不需要手工把 `CERT_DOMAIN` 移到第三列，也不需要额外的 `reality-check-csv` 包装脚本。文件名和目录可以任意，只要当前用户有读取权限。

## 数据文件

首次运行时，程序会检查 `data/` 目录中的 GeoIP、GFWList、CDN 关键词和热门网站数据；缺少时会尝试自动下载。不要把本地 CSV、扫描结果、数据库、缓存文件或编译后的二进制提交到 GitHub。

## 从其他机器检测

可以在本地运行 RealiTLScanner，把 CSV 上传到安装了 RealityChecker 的机器：

```bash
scp file2.csv user@server:/path/to/workdir/
ssh user@server
./reality-checker csv /path/to/workdir/file2.csv
```

也可以在另一台机器重新下载同一个 Release 或重新编译源码。

## Windows / macOS

```bash
go build -o reality-checker .
./reality-checker csv file.csv
```

Windows PowerShell：

```powershell
go build -o reality-checker.exe .
.\reality-checker.exe csv .\file.csv
```

## 开发与测试

```bash
go test ./...
gofmt -w internal/cmd/csv.go internal/cmd/csv_test.go
```

## 免责声明

本工具仅用于技术研究和学习，请遵守当地法律法规以及目标网络的使用政策。
