# Distributed Task Scheduler

[![Go Version](https://img.shields.io/badge/go-1.21%2B-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://img.shields.io/github/actions/workflow/status/yourusername/taskscheduler/go.yml?branch=main)](https://github.com/yourusername/taskscheduler/actions)

A high-performance, distributed task scheduler written in Go, designed to handle millions of scheduled tasks with built-in retry mechanisms, persistence, and observability.

## Features

- **Millions of Tasks** - Efficient priority queue implementation using min-heap
- **Retry Mechanisms** - Configurable retry policies with exponential backoff
- **Persistence** - Redis-backed task storage for fault tolerance
- **Observability** - Prometheus metrics and Grafana dashboard support
- **Rate Limiting** - Token bucket rate limiting for task execution
- **Priority Scheduling** - Support for task prioritization
- **Docker Support** - Containerized deployment with Docker Compose
- **Worker Pool** - Configurable concurrent workers
- **Graceful Shutdown** - Proper cleanup of resources on termination

## Build and Run

### Prerequisites
- Go 1.21+
- Redis 6.2+
- Docker 20.10+ (optional)

### Build from Source

# Clone repository
git clone https://github.com/ichbingautam/distributed-task-scheduler.git
cd distributed-task-scheduler

# Install dependencies
go mod tidy

# Build binary
go build -o scheduler ./cmd/scheduler/main.go

# Run with default configuration
./scheduler

## Architecture

```mermaid
graph TD
    A[Client] --> B(Scheduler API)
    B --> C[Task Queue]
    C --> D[Redis Storage]
    D --> E[Executor]
    E --> F[(Metrics)]
    E --> G[Worker Pool]
    G --> H[(External Services)]