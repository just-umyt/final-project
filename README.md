# 🎓 Final Project – Backend Development Bootcamp

This is my final project for the **Backend Development Bootcamp** at **Baky Tylla**.  
It consists of multiple microservices written in Go, integrated with modern backend technologies.

---

## 🏗️ General Architecture

![General service architecture](docs/img/General%20Project%20Architecture.png)

This project is composed of the following services:

- **🛒 Cart Service** (`cart`) – Manages user shopping cart operations.
- **📦 Stocks Service** (`stocks`) – Manages product inventory and prices.
- **🌀 Kafka** (`kafka`) – Handles message brokering, ZooKeeper, and cluster configuration.
- **📊 Metrics Consumer** (`metrics-consumer`) – Consumes Kafka events and logs them.
- **🧪 Monitoring** (`monitoring`) – Provides logging, tracing, and metrics with Prometheus, Grafana, and Jaeger.

Each service has its own documentation and instructions on how it works and how to test it.  
_📁 Note: You’ll also find a `proto/` folder used for gRPC – no need to focus on it._

---

## ⚙️ Technologies Used

This project leverages a wide range of technologies:

- **Programming**: Go (Golang)
- **Communication**: gRPC
- **Database**: PostgreSQL
- **Message Broker**: Apache Kafka
- **Containerization**: Docker, Docker Compose
- **Monitoring & Observability**:
  - Prometheus (metrics)
  - Grafana (visualization)
  - Jaeger (distributed tracing)
  - OpenTelemetry
- **Logging**: Zap
- **Testing**: Minimock (for mocking)

📂 Project Structure:

```

.
├── cart/
├── stock/
├── kafka/
├── metrics-consumer/
├── monitoring/
└── proto/

```

📊 Technologies overview:

![All technologies](docs/img/All%20techs.png)

---

## 🧪 How to Run the Project

1. Create a Docker network (only once):

```bash
docker network create public-net
```

2. Start all services using `make`:

```bash
make docker-up
```

3. Access services on the following ports:

| Service          | Description                   | Port    |
| ---------------- | ----------------------------- | ------- |
| 🛒 Cart Service  | HTTP Gateway                  | `8080`  |
| 📦 Stock Service | HTTP Gateway                  | `8081`  |
| 📈 Prometheus    | Monitoring (metrics)          | `9090`  |
| 📉 Grafana       | Dashboards & visualization    | `3000`  |
| 🧭 Jaeger UI     | Distributed tracing interface | `16686` |

---

## Monitoring

## 🙏 Acknowledgments

I would like to express my sincere gratitude to **Baky Tylla Education Center** and the entire team for their guidance and support during the Backend Development Bootcamp.
Thank you for helping me learn and apply these powerful backend technologies!
