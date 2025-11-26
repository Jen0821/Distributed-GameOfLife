## 🚀 Go-Distributed-GameOfLife

This project implements **Conway's Game of Life** as a **highly scalable, distributed system** using **Go (Golang)**. The primary focus is on achieving **high-performance simulation** across multiple machines using a **Controller/Server/Worker** architecture and Go's native support for network communication via **RPC (Remote Procedure Call)**.

This system is designed to distribute the computational load of large-scale grids across several independent nodes (simulated **AWS Nodes**), demonstrating exceptional horizontal scalability.

## Overview

Conway's Game of Life is a zero-player game whose evolution is determined solely by its initial state. The grid of cells evolves based on simple, local rules:

1.  Any live cell with fewer than two live neighbours dies, as if by underpopulation.
2.  Any live cell with two or three live neighbours lives on to the next generation.
3.  Any live cell with more than three live neighbours dies, as if by overpopulation.
4.  Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.

This project's implementation focuses on distributing these intense computations across multiple remote nodes using the Controller/Server/Worker model to break through the computing power bottleneck of a single machine.

## 📝 Final Coursework Report

The full detailed analysis of the distributed implementation, including performance benchmarks and design rationale, is available in the final report.

[![Report Preview](https://github.com/Jen0821/Distributed-GameOfLife/blob/main/report.jpg)](https://github.com/Jen0821/Distributed-GameOfLife/blob/main/report.pdf)

**Click the image to view the full report.**

## 💡 System Architecture: Controller-Server-Worker

![Architecture Diagram](architecture.png)

The simulation logic is split into three main, network-connected components—the **Local Controller** and the **Remote Engine (Server/Worker)**—to maximize performance and scalability.

| Component | Role in Distributed System | Key Technologies |
| :--- | :--- | :--- |
| **Local Controller (Distributor/Client)** | **Orchestrates the system**, handles user input/visualization (SDL), initiates **RPC calls** to remote Workers, collects calculation results, and **aggregates them** into the global world state. | Go, SDL, RPC Client |
| **Engine (Remote Server)** | Deployed on **AWS nodes**. Receives RPC requests via **`RunServer`**, holds the core GOL calculation logic (`DoWork`), and returns the results. | Go, RPC Server/Client, Goroutines |
| **GOL Worker (Compute Logic)** | The calculation logic within the Engine. It calculates the GOL logic for its assigned board **slice (shard)** using **Goroutines** for internal parallelism. | Go, Goroutines (internal parallelism) |

> **Note on Broker Role:** In the current implementation, the **Local Controller** handles the scheduling and result aggregation roles. A dedicated Broker component for managing worker pools and fault tolerance was part of the original design but is functionally handled by the Controller.

## 🌟 Key Features & Implementation Highlights

### 1\. High-Performance Distributed Core

  * **RPC Communication:** All communication between the Controller and the Server is handled via **RPC (Remote Procedure Call)**, ensuring reliable command and data transfer over the network.
  * **Workload Partitioning:** The global board is divided evenly by `ImageHeight` into multiple regions (**shards**), and each remote Server/Worker calculates its assigned row shard.
  * **Halo Data Transmission:** To correctly compute boundary cells, the **Halo (boundary data)** (`HaloUpper`, `HaloLower`) is explicitly defined and sent from the Controller to the Workers with each request. This prevents complex neighbour communication between remote Workers, simplifying the design.
  * **RPC Structures:** Communication relies on defined RPC structures:
      * **`WorkerRequest`:** Carries initial parameters (e.g., `StartY`, `Height`, `World`, `HaloUpper`, `HaloLower`) from the Controller to the Worker.
      * **`WorkerResponse`:** Holds the calculated results from the Worker, returning the calculated area status through `Result[][]byte`.

### 2\. Solved Implementation Challenges

The core distributed architecture was successfully established within the **`StartDistributed{}`** function in the distributor:

  * **Controller-Server Connection:** The Controller successfully connects to remote Servers using `rpc.Dial` and executes remote functions via `client.call()`.
  * **Boundary Cell Dependency:** The issue of boundary cells missing neighbouring rows was solved by embedding the **`HaloUpper`** and **`HaloLower`** data within the `WorkerRequest` for every turn's calculation.

### 3\. State Management & Real-Time Events

  * **Toroidal Boundary Conditions:** Implements **Closed Domain (Toroidal)** boundary conditions.
  * **Event-Driven Reporting:** The Controller receives the aggregated number of surviving cells from the Workers for real-time notification.
  * **PGM I/O:** Uses the **PGM (Portable Graymap)** image format for loading initial states and saving final/intermediate states.

### 4\. Interactive User Control

Interactive keyboard commands are processed by the Local Controller and transmitted as RPC commands to the remote Server/Engine:

  * **`s` (Save):** Commands the remote system to save the current simulation state as a PGM image.
  * **`q` (Quit):** Gracefully terminates the system and triggers the final PGM image output.
  * **`p` (Pause/Resume):** Toggles the processing state on the remote Server/Worker nodes.

## 📈 Performance and Scalability

The distributed architecture is optimized for **horizontal scaling**. Benchmarks were conducted on a 20-core Linux machine using a $512 \times 512$ grid over 1000 turns.

  * **Network Overhead:** Initially, running time **significantly increased** compared to the single-machine parallel implementation due to the time loss caused by calling remote Workers (network latency).
  * **Scaling Proof:** As the number of Workers increased, the running time for completing the Life Game **significantly decreased**, confirming the effectiveness of horizontal distribution.
  * **Stabilisation:** As the number of Workers continued to increase, the improvement in running time became smaller and eventually **stabilised** due to increasing network and aggregation overhead.

The chart below illustrates the speedup achieved by distributing the workload:

| Workers Number | Runtime (s) (Approx. from Chart) |
| :--- | :--- |
| 1 | \~210 |
| 2 | \~120 |
| 3 | \~90 |
| 4 | \~70 |
| 5 | \~55 |
| 6 | \~45 |

### Potential Faults

The design considers potential fault scenarios inherent to distributed systems:

  * **Network Interruption:** If the network connection is interrupted, the RPC requests or responses will fail, making the system unable to continue calculating and aggregating the next state.
  * **Worker Anomaly:** If a Worker experiences anomalies, the local server will lack the states of certain shards, resulting in errors in the aggregation of the global GOL state.
  * **Controller Crash:** If the local server (Controller) suddenly crashes, the current world state cannot be saved, and recovery may require starting from the initial state.

## ▶️ Setup

### **Prerequisites**

Install Go and SDL development libraries.

### macOS (Homebrew)

```bash
brew install sdl2
```

### Linux (Ubuntu/Debian)

```bash
sudo apt-get install libsdl2-dev
```

### Windows

```bash
# Requires MinGW installation and manual SDL2 linking
# Refer to the Go SDL documentation for detailed platform-specific setup
```

## 🚀 Running and Testing

### Run the program (Controller/Server assumed running or mocked)

```bash
go run .
```

### Test core functionality (Step 1 and 2 tests passed smoothly)

```bash
go test tests -v -run TestGol/-1\$
```

### Run tests with race detector

```bash
go test ./tests -v -race
```
