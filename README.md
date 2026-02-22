# 🕊️ Salam

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go&logoColor=white)](https://github.com/theb0imanuu/salam/blob/main/go.mod)
[![Go Report Card](https://goreportcard.com/badge/github.com/theb0imanuu/salam)](https://goreportcard.com/report/github.com/theb0imanuu/salam)
[![License](https://img.shields.io/github/license/theb0imanuu/salam)](https://github.com/theb0imanuu/salam/blob/main/LICENSE)
[![Repo Size](https://img.shields.io/github/repo-size/theb0imanuu/salam)](https://github.com/theb0imanuu/salam)
[![GitHub Stars](https://img.shields.io/github/stars/theb0imanuu/salam?style=social)](https://github.com/theb0imanuu/salam/stargazers)
[![GitHub Issues](https://img.shields.io/github/issues/theb0imanuu/salam)](https://github.com/theb0imanuu/salam/issues)
[![GitHub Forks](https://img.shields.io/github/forks/theb0imanuu/salam)](https://github.com/theb0imanuu/salam/network/members)

**Salam** (سلام) means "peace" in Arabic. A fast, lightweight server health monitoring CLI written in Go.

## Features

- ⚡ **Fast**: Compiled Go binary, no runtime dependencies
- 📊 **Interactive Dashboard**: Real-time TUI dashboard with live updates
- 🎯 **Precise**: Detailed CPU, memory, disk, and network metrics
- 🔔 **Alerts**: Configurable thresholds with webhook support
- 📊 **Multiple Formats**: Pretty console output or JSON
- 🖥️ **Cross-Platform**: Linux, macOS, Windows support
- 🔧 **Easy Install**: Single binary, global installation

## Setup Instructions

1. **Initialize the project:**

   ```bash
   mkdir salam
   cd salam
   go mod init github.com/theb0imanuu/salam
   ```

2. **Create the directory structure:**

   ```bash
   mkdir -p cmd/salam internal/{monitor,reporter,config,models} pkg/utils
   ```

3. **Copy the files into their respective locations**

4. **Download dependencies:**

   ```bash
   go mod tidy
   ```

5. **Build:**

   ```bash
   make build
   ```

6. **Test locally:**

   ```bash
   ./build/salam check
   ```

7. **Install globally:**

   ```bash
   make install
   # or
   cp build/salam /usr/local/bin/
   ```

8. **Create releases:**
   ```bash
   make release
   make package
   ```

## Installation

### Via Installer (Recommended)

Download the latest binary for your OS from [releases](https://github.com/theb0imanuu/salam/releases).

1. **Windows**: Run `salam-setup-windows-amd64.exe`. It will install Salam and add it to your PATH automatically.
2. **Unix**: Run the `salam-install` binary.

### Via Install Script (Unix only)

### Via Go

```bash
go install github.com/theb0imanuu/salam/cmd/salam@latest
```

### Manual Download

Download the latest binary from [releases](https://github.com/theb0imanuu/salam/releases).

## Usage

### Interactive Dashboard (v2.0 preview)

```bash
salam dashboard
```

### Quick Check

```bash
salam check
```

### Check Specific Metrics

```bash
salam check --cpu
salam check --memory
salam check --disk
salam check --network
```

### Continuous Monitoring

```bash
salam watch --interval 30 --threshold 80
```

### JSON Output

```bash
salam check --json > report.json
```

### Webhook Alerts

```bash
salam watch --webhook https://hooks.slack.com/services/...
```

## Configuration

Generate a config file:

```bash
salam config
```

Edit `~/.salam.yaml`:

```yaml
thresholds:
  cpu: 80
  memory: 85
  disk: 90
  load: 2.0

alerts:
  enabled: true
  webhook: "https://your-webhook-url.com"
  email: "admin@example.com"
```

## Building from Source

```bash
git clone https://github.com/theb0imanuu/salam.git
cd salam
make build
```

## License

MIT License - see [LICENSE](file:///c:/Users/MANU/Desktop/salam/LICENSE) file for details.
