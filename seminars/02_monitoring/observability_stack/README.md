# Go Web Services Observability Stack

This setup uses modern, high-performance alternatives to traditional tools.

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose
- Go

### 1. Setup Observability Stack
```bash
# Start all observability services
./setup.sh
```

### 2. Run Go Application Locally

### 3. Access the Services
- **Grafana**: http://localhost:3000 (admin/admin)
- **Victoria Metrics**: http://localhost:8428"
- **Victoria Logs**: http://localhost:9428"
- **Go Application**: http://localhost:8080"
- **Node Exporter**: http://localhost:9100"

### 4. Stop Everything
```bash
# Stop the Go application (Ctrl+C)
# Stop observability services
docker-compose down
```

## 📊 Observability Stack Components

### 📈 Metrics & Storage
- **Victoria Metrics** - High-performance metrics storage (faster than Prometheus)
  - URL: http://localhost:8428
  - Web UI available at the above URL

### 📝 Log Aggregation
- **Victoria Logs** - Fast log aggregation and querying (faster than Loki)
  - URL: http://localhost:9428
  - Query logs using LogQL syntax

### 🎨 Visualization
- **Grafana** - Data visualization and dashboarding
  - URL: http://localhost:3000
  - Username: `admin`
  - Password: `admin`
  - Pre-configured dashboards included

### 🐛 Error Tracking
- **Sentry** - Error monitoring and performance tracking
  - URL: http://localhost:9000
  - Setup required on first visit

### 📡 Telemetry Collection
- **Victoria Metrics Agent (vmagent)** - Metrics collection and forwarding
  - Scrapes metrics from the Go application
  - Forwards metrics to Victoria Metrics

### 🖥️ System Monitoring
- **Node Exporter** - System metrics collection
  - URL: http://localhost:9100/metrics

## 📚 Topics

This setup demonstrates key observability concepts:

### 1. **Metrics**
- Request rate, duration, and error rates
- Go runtime metrics (GC, memory, goroutines)
- System metrics (CPU, memory, disk)

### 2. **Logging**
- Structured logging with correlation IDs
- Log aggregation and searching
- Log levels and filtering

### 3. **Tracing**
- Distributed tracing concepts
- Span creation and propagation
- Trace sampling strategies

### 4. **Error Tracking**
- Automatic error capture
- Stack trace analysis
- Error grouping and alerting

## 🔧 Configuration

### Environment Variables

The Go application should use these environment variables:

```bash
# Metrics Configuration (for Victoria Metrics)
VM_EXPORTER_URL=http://vmagent:8429

# Sentry Configuration
SENTRY_DSN=YOUR_SENTRY_DSN_HERE
SENTRY_ENVIRONMENT=development
```

### Victoria Metrics Agent Configuration

The agent is configured in `vmagent-config.yml` with **Scraping**: Direct scraping from Go application and other services

### Grafana Data Sources

Pre-configured data sources:
- Victoria Metrics (primary metrics)
- Victoria Logs (log aggregation)
- etc... todo: add preconf

## 📖 Key Concepts

### The Three Pillars of Observability
1. **Metrics**: Quantitative data about system behavior
2. **Logs**: Detailed event records with context
3. **Traces**: Request flow through distributed systems

### Go-Specific Observability
- Runtime metrics (GC, goroutines, memory)
- Built-in profiling support
- Structured logging best practices

## 🚦 Service URLs

| Service          | URL                    | Purpose             |
|------------------|------------------------|---------------------|
| Go Application   | http://localhost:8080  | Main application    |
| Grafana          | http://localhost:3000  | Visualization       |
| Victoria Metrics | http://localhost:8428  | Metrics storage     |
| Victoria Logs    | http://localhost:9428  | Log aggregation     |
| Node Exporter    | http://localhost:9100  | System metrics      |

## 🔍 Troubleshooting

### Common Issues

1. **Services not starting**
   ```bash
   # Check Docker logs
   docker-compose logs [service-name]
   
   # Check port conflicts
   netstat -tulpn | grep [port]
   ```

2. **No metrics in Grafana**
   - Verify vmagent is running and scraping metrics
   - Check Go application environment variables
   - Verify data source configuration

3. **Missing logs in Victoria Logs**
   - Check vmagent log pipeline
   - Verify log format and labels
   - Check retention settings

4. **Sentry not capturing errors**
   - Update SENTRY_DSN environment variable
   - Verify Sentry service is running
   - Check network connectivity

## 📚 Additional Resources

- [Victoria Metrics Documentation](https://docs.victoriametrics.com/)
- [Grafana Documentation](https://grafana.com/docs/)
- [Sentry Documentation](https://docs.sentry.io/)
