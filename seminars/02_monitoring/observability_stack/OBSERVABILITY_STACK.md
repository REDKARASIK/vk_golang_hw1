# Go Web Services Observability Stack

## 🏗️ Architecture Overview

The observability stack follows the three pillars of observability:
- **Metrics**: Victoria Metrics (modern Prometheus alternative)
- **Logs**: Victoria Logs (modern Loki alternative)  

Additional tools:
- **Visualization**: Grafana with pre-configured dashboards todo: add preconf
- **Error Tracking**: Sentry for bug monitoring
- **System Monitoring**: Node Exporter
- **Data Collection**: OpenTelemetry Collector

## 📦 Components

### Core Observability Tools

| Tool                 | Purpose                       | Port  | Why This Choice                              |
|----------------------|-------------------------------|-------|----------------------------------------------|
| **Victoria Metrics** | Metrics storage and querying  | 8428  | Faster than Prometheus, lower memory usage   |
| **Victoria Logs**    | Log aggregation and searching | 9428  | More efficient than Loki, better performance |
| **Grafana**          | Visualization and dashboards  | 3000  | Most popular visualization tool              |

### Infrastructure Monitoring

| Tool                        | Purpose                   | Port      |
|-----------------------------|---------------------------|-----------|
| **Node Exporter**           | System metrics            | 9100      |

## 🚀 Quick Start

### 1. Setup
```bash
# Validate configuration
./validate-config.sh

# Start the stack
./setup.sh
```

### 2. Access URLs

| Service          | URL                    | Credentials |
|------------------|------------------------|-------------|
| Grafana          | http://localhost:3000  | admin/admin |
| Victoria Metrics | http://localhost:8428  | -           |
| Victoria Logs    | http://localhost:9428  | -           |

### 3. Integration with Go Applications

#### Metrics Integration
```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "net/http"
)

var (
    requestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
)

func init() {
    prometheus.MustRegister(requestsTotal)
    http.Handle("/metrics", promhttp.Handler())
}
```

#### Logging Integration
```go
import (
    "log/slog"
    "os"
)

func initLogger() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    logger.Info("time in handler", "name", name, "time", "123ms")
}
```

## 🔧 Configuration Files

### docker-compose.yaml
Main orchestration file defining all services, networks, and volumes.

### vmagent-config.yml
Victoria Metrics Collector (vmagent) configuration for receiving, processing, and exporting metrics.

### Grafana Provisioning
- **Datasources**: Automatic configuration of Victoria Metrics and Victoria Logs

### Best Practices
1. Use structured logging ("log/slog")
2. Add meaningful attributes to spans
3. Create custom metrics for business logic
4. Set up appropriate alerting thresholds
5. Regularly review and optimize dashboards

## 🛑 Cleanup

To stop and remove all services:
```bash
docker-compose down -v
```

To remove all data (be careful!):
```bash
docker-compose down -v --remove-orphans
docker system prune -f
```
