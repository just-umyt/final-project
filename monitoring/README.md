# 📊 Logging, Tracing, and Metrics

## 🧩 Overview

In this assignment, we focused on enhancing **observability** in our microservices by implementing the following:

1. **Structured Logging** – for consistent and informative logs.
2. **Distributed Tracing** – for tracking requests across services.
3. **Prometheus Metrics** – for collecting and monitoring performance data.

---

## 🎯 Objectives

### 📝 Logging

- Used a structured logging library (e.g., **Zap**) to record meaningful events.
- Included contextual information such as:
  - gRPC method and path
  - Trace ID
  - Error messages (if any)

### 📍 Tracing

- Integrated **OpenTelemetry** for tracing.
- Exported traces to **Jaeger** for visualization.
- Propagated trace context between services using HTTP/gRPC headers.
- Each request includes:
  - A unique **Trace ID**
  - One or more **Spans** to represent operations

### 📈 Metrics

Collected using **Prometheus**:

| Metric Name             | Type      | Description                             |
| ----------------------- | --------- | --------------------------------------- |
| `failed_requests_total` | Counter   | Increments on every failed HTTP request |
| `response_time_seconds` | Histogram | Measures the duration of each HTTP call |

---

## ⚙️ Quick Tip for Grafana

> 💡 You can quickly get started in Grafana using the **default `Grafana default dashboards`** template.
