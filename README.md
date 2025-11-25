# 🚀 Go-Distributed-GameOfLife

This project implements **Conway's Game of Life** as a **highly scalable, distributed system** using **Go (Golang)**. [cite\_start]The primary focus is on achieving **high-performance simulation** across multiple machines using a **Controller/Server/Worker** architecture and Go's native support for network communication via **RPC (Remote Procedure Call)**[cite: 113, 114, 128].

[cite\_start]This system is designed to distribute the computational load of large-scale grids across several independent nodes (simulated **AWS Nodes**)[cite: 115, 127], demonstrating exceptional horizontal scalability.

## Overview

Conway's Game of Life is a zero-player game whose evolution is determined solely by its initial state. The grid of cells evolves based on simple, local rules:

1.  Any live cell with fewer than two live neighbours dies, as if by underpopulation.
2.  Any live cell with two or three live neighbours lives on to the next generation.
3.  Any live cell with more than three live neighbours dies, as if by overpopulation.
4.  Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.

[cite\_start]This project's implementation focuses on distributing these intense computations across multiple remote nodes using the Controller/Server/Worker model to break through the computing power bottleneck of a single machine[cite: 127].

-----

## 📝 Final Coursework Report

The full detailed analysis of the distributed implementation, including performance benchmarks and design rationale, is available in the final report.

[](https://www.google.com/search?q=./report.pdf)

### Direct Link

[Download the Full PDF Report Here](https://www.google.com/search?q=./report.pdf)

-----

## 💡 System Architecture: Controller-Server-Worker

[cite\_start]The simulation logic is split into two primary, network-connected parts—the **Local Controller** and the **Remote Engine (Server)**—to maximize performance and scalability[cite: 126].

| Component | Role in Distributed System | Key Technologies |
| :--- | :--- | :--- |
| **Local Controller (Distributor/Client)** | [cite\_start]Orchestrates the system, handles user input/visualization (SDL), initiates **RPC calls** to remote Workers, collects calculation results, and merges them into the global world state[cite: 125, 130, 139]. | Go, SDL, RPC Client |
| **Engine (Remote Server)** | [cite\_start]Deployed on **AWS nodes**[cite: 115, 127]. [cite\_start]Receives RPC requests, holds the core GOL calculation functionality (`DoWork{}`) [cite: 119, 120][cite\_start], and returns the results[cite: 130]. | [cite\_start]Go, RPC Server/Client, `RunServer` [cite: 129] |
| **GOL Worker (Compute Logic)** | The calculation logic within the Engine. [cite\_start]It calculates the GOL logic for its assigned board **slice (shard)**[cite: 150]. | Go, Goroutines (internal parallelism) |

> [cite\_start]**Note on Broker:** The construction of an ideal **Broker** component (managing node status, allocating tasks to idle Workers, and aggregating results) was attempted but not completed[cite: 164, 166, 167, 168]. The **Local Controller** currently handles the scheduling and aggregation roles.

-----

## 🌟 Key Features & Implementation Highlights

### 1\. High-Performance Distributed Core

  * [cite\_start]**RPC Communication:** All communication between the Controller and the Server is handled via **RPC (Remote Procedure Call)** [cite: 128][cite\_start], which enables multi-machine collaboration via network transmission[cite: 131].
  * [cite\_start]**Workload Partitioning:** The world is divided evenly by `ImageHeight` into multiple regions (**shards**), and each Worker solves the calculation for its assigned row shard[cite: 149, 150]. [cite\_start]This simplifies parameter transmission[cite: 151].
  * [cite\_start]**Halo Data Transmission:** The **Halo (boundary data)** (`HaloUpper`, `HaloLower`) is explicitly defined and sent from the Controller to the Workers with each request to solve the problem of boundary cells missing neighboring rows from adjacent shards[cite: 181, 182, 183].
  * **RPC Structures:** Communication uses defined RPC structures:
      * [cite\_start]**`WorkerRequest`:** Carries initial parameters from the Controller to the Worker (e.g., `StartY`, `Height`, `World`, `HaloUpper`, `HaloLower`)[cite: 155, 156, 122].
      * [cite\_start]**`WorkerResponse`:** Holds the calculation results from the Worker, returning the calculated area status through `Result[][]byte`[cite: 122, 158].

### 2\. Solved Implementation Challenges

The following challenges were overcome during development:

  * [cite\_start]**Controller-Server Connection:** The core distributed implementation was built within the **`StartDistributed{}`** function in the distributor, allowing the Controller to successfully connect to remote Workers using `rpc.Dial` and `client.call()`[cite: 176, 178].
  * [cite\_start]**Boundary Cell Dependency:** The issue of boundary cells missing neighboring rows was solved by transmitting the **`HaloUpper`** and **`HaloLower`** data with the initial shard information for every turn's calculation[cite: 181, 182, 183].

### 3\. State Management & Real-Time Events

  * **Toroidal Boundary Conditions:** Implements **Closed Domain (Toroidal)** boundary conditions.
  * [cite\_start]**Event-Driven Reporting:** The Controller receives the number of surviving cells from the Worker for event notification[cite: 161].
  * **PGM I/O:** Uses the **PGM (Portable Graymap)** image format for loading initial states and saving final/intermediate states.

### 4\. Interactive User Control

[cite\_start]Interactive keyboard commands are processed by the Local Controller and transmitted as RPC commands to the remote Server/Engine[cite: 125, 162]:

  * **`s` (Save):** Commands the remote system to save the current simulation state as a PGM image.
  * **`q` (Quit):** Gracefully terminates the system and triggers the final PGM image output.
  * **`p` (Pause/Resume):** Toggles the processing state on the remote Worker nodes.

-----

## 📈 Performance and Scalability

The distributed architecture is optimized for **horizontal scaling**. [cite\_start]Benchmarks were conducted on a 20-core Linux machine [cite: 184] using a $512 \times 512$ grid over 1000 turns.

  * [cite\_start]**Network Overhead:** Initially, the running time **significantly increased** compared to the parallel implementation due to the time loss caused by calling remote Workers (network latency)[cite: 190, 191].
  * [cite\_start]**Scaling Proof:** As the number of Workers increased, the running time for completing the Life Game **significantly decreased**[cite: 192].
  * [cite\_start]**Stabilization:** As the number of Workers continued to increase, the improvement in running time became smaller and eventually stabilized[cite: 193].

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

[cite\_start]The design acknowledges potential fault scenarios that can occur in a distributed system[cite: 209]:

  * [cite\_start]**Network Interruption:** If the network connection is interrupted, the RPC requests or responses will fail, making the system unable to continue calculating and aggregating the next state[cite: 210, 211].
  * [cite\_start]**Worker Anomaly:** If a Worker experiences an anomaly, the local server will lack the states of certain shards, resulting in errors in the aggregation of the global GOL state[cite: 212].
  * [cite\_start]**Controller Crash:** If the local server (Controller) crashes, the current world state cannot be saved, and recovery may require starting from the initial state[cite: 213].

-----

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
