# 🚀 Go-Distributed-GameOfLife

This project implements **Conway's Game of Life** as a **highly scalable, distributed system** using **Go (Golang)**. The primary focus is on achieving **high-performance simulation** across multiple machines using a **Controller/Broker/Worker** architecture and Go's native support for network communication via **RPC (Remote Procedure Call)**.

This system is designed to distribute the computational load of large-scale grids across several independent nodes (simulated **AWS Nodes**), demonstrating exceptional horizontal scalability.

## Overview

Conway's Game of Life is a zero-player game whose evolution is determined solely by its initial state. The grid of cells evolves based on simple, local rules:

1.  Any live cell with fewer than two live neighbours dies, as if by underpopulation.
2.  Any live cell with two or three live neighbours lives on to the next generation.
3.  Any live cell with more than three live neighbours dies, as if by overpopulation.
4.  Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.

This project's implementation focuses on distributing these intense computations across multiple remote nodes using the Controller/Broker/Worker model to achieve massive scalability.

## 📝 Final Coursework Report

The full detailed analysis of the parallel implementation, including performance benchmarks and design rationale, is available in the final report.

[![Report Cover Image: Click to Download PDF](./report.jpg)](./report.pdf)
### Direct Link
[Download the Full PDF Report Here](./report.pdf)

## 💡 System Architecture: Controller-Broker-Worker

The simulation logic is split into three main, network-connected components to maximize performance and fault tolerance potential.

![Architecture Diagram](architecture.png)

| Component | Role in Distributed System | Key Technologies |
| :--- | :--- | :--- |
| **Local Controller (Client/UI)** | Responsible for user input, **SDL** real-time visualization, sending RPC commands, and receiving events from the Broker. | Go, SDL, RPC Client |
| **Broker (Server/Orchestrator)** | Central management point. Manages the Worker node pool, handles board slicing/distribution, and aggregates results (Aggregation) from Workers. | Go, RPC Server/Client |
| **GOL Worker (Compute Node)** | Calculates the GOL logic for its assigned board slice. Uses **Goroutines** for internal parallelism and exchanges boundary data (**Halo**) with neighboring Workers. | Go, Goroutines, RPC Client |

## 🌟 Key Features & Implementation Highlights

### 1. High-Performance Distributed Core
* **RPC Communication:** All communication between components (Controller, Broker, Worker) is handled via **RPC (Remote Procedure Call)**, ensuring reliable command and data transfer over the network.
* **Efficient Workload Partitioning:** The Broker slices the global board horizontally and assigns sections to Worker nodes. **Halo Exchange** is used to efficiently communicate only the boundary cell information between neighbors, minimizing network overhead.
* **Internal Parallelism:** Each Worker node utilizes Go's **Goroutines** internally for fast, local parallel computation on its assigned board slice.

### 2. State Management & Real-Time Events
* **Toroidal Boundary Conditions:** Implements **Closed Domain (Toroidal)** boundary conditions, where opposite edges of the board are connected.
* **Event-Driven Reporting:**
    * **Alive Count Ticker:** Reports the total number of alive cells to the Controller every **2 seconds** via an RPC call.
    * **State Change Events:** Utilizes `CellFlipped` and `TurnComplete` events for efficient real-time visualization updates.
* **PGM I/O:** Uses the **PGM (Portable Graymap)** image format for loading initial states and saving final/intermediate states.

### 3. Interactive User Control
Interactive keyboard commands are processed by the Local Controller and sent as RPC commands to the distributed system:
* **`s` (Save):** Commands the remote system to save the current simulation state as a PGM image.
* **`q` (Quit):** Gracefully terminates the Controller client and triggers the final PGM image output.
* **`p` (Pause/Resume):** Toggles the processing state on the remote Worker nodes.

## 📈 Performance and Scalability

This architecture is optimized for **horizontal scaling**. The final report's benchmarks illustrate key performance characteristics:

* **Scalability Proof:** The time required to complete a fixed number of turns **significantly decreases** as the number of Worker nodes increases, demonstrating the effectiveness of the distributed approach.
* **Overhead Analysis:** The performance gain eventually plateaus due to the increasing overhead of network communication and Broker aggregation.
* **Fault Consideration:** The design considers potential fault scenarios, such as handling a new Controller taking over the session or a Worker Node failing.

## ▶️ Setup and Running

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

## 🚀 Running

### Run the program (Controller/Broker/Workers assumed running or mocked)
```bash
go run .
```

### Test visualization + keyboard controls
```bash
go test ./tests -v -run TestKeyboard -sdl
```

### Test parallel core with race detector
```bash
go test ./tests -v -race
```
