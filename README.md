# 🌐 Go-Distributed-GameOfLife: Scalable GOL Engine (Conway's Game of Life distributed Implementation)

This repository hosts a high-performance, **distributed implementation** of **Conway's Game of Life (GoL)**. Built entirely in **Go (Golang)**, the system leverages an asynchronous, network-based architecture using **Remote Procedure Calls (RPC)** to distribute computation across multiple machines (simulated as AWS nodes), achieving massive scalability for large grid simulations.

The core design separates I/O and visualization (Local Controller) from the primary computational load (AWS Nodes/Workers).

## 🔍 Overview

**Distributed-GameOfLife-Go** is a fully parallel, distributed implementation of  
**Conway’s Game of Life**, designed for high-performance simulations running across
multiple machines or cloud nodes (e.g., AWS EC2).

The system is architected around a **three-layer distributed model**:

- **Local Controller** — handles visualization, PGM I/O, keyboard events  
- **Broker / Distributor** — coordinates work, splits the board into slices  
- **Worker Nodes** — compute next-generation slices using halo exchange  

This separation ensures:

- Extremely scalable performance  
- Modular component replacement  
- Low coupling between visualization and computation  
- High fault tolerance (workers can join/leave dynamically)

The architecture enables real-time visualization, efficient PGM save/load operations,
and clean RPC-based communication between all components.

## 📝 Final Coursework Report

The full detailed analysis of the parallel implementation, including performance benchmarks and design rationale, is available in the final report.

[![Report Cover Image: Click to Download PDF](./report.jpg)](./report.pdf)
### Direct Link
[Download the Full PDF Report Here](./report.pdf)

## ✨ Distributed System Architecture

The project follows a three-tier distributed model designed for clarity, scalability, and workload separation:

1. **Local Controller (Client):**  
   Handles user I/O (PGM load/save, SDL visualization, and keypresses).  
   Sends commands to the **Broker** via RPC.

2. **Broker / Distributor (Central Server):**  
   Acts as the system orchestrator.  
   Receives commands from the Controller, splits the global board into slices, and sends work to Workers.  
   Contains **no game logic**.

3. **GOL Worker (Compute Node):**  
   Runs on separate machines (AWS nodes / terminals).  
   Computes the next state for a given board slice.  
   Uses **Halo Exchange** to safely exchange boundary rows with neighboring Workers.

### 📊 Architecture Diagram

![Distributed GOL Architecture](./architecture.png)

## ⚙️ Running the Distributed System

This system requires launching the components in the following order: **Broker → Workers → Controller**.

### Prerequisites

- Go (Golang) installed  
- SDL2 development libraries (only for visualization on the Controller)

## 🔌 Execution Sequence

### 1. Start the Broker / Distributor

The Broker listens for RPC calls from both the Controller and Worker nodes.

```bash
# Start the Broker on a specified port (e.g., 8080)
go run ./broker -listen :8080
```

### 2. Start the GOL Worker Nodes

Each Worker connects to the Broker and waits for assigned board slices.  
Start as many Workers as you want (more Workers = more parallelism).

```bash
# Start Worker 1
go run ./worker -broker-addr :8080 -worker-id 1

# Start Worker 2 (different machine/terminal)
go run ./worker -broker-addr :8080 -worker-id 2

# ... and so on for N workers
```

### 3. Start the Local Controller (Client / I/O)

The Local Controller manages SDL display, PGM save/load, and user key input.  
Use the `-headless` flag for performance testing without rendering.

```bash
# Start the Controller connected to the Broker
go run ./controller -broker-addr :8080 -headless -t 4
```

## 🕹️ Interactive Controls (Controller)

| Key | Action | Description |
|-----|--------|-------------|
| **s** | Save State | Broker assembles the global board and writes a PGM image |
| **p** | Pause / Resume | Toggles simulation updates while preserving control commands |
| **q** | Quit Client | Controller exits gracefully and outputs final PGM |
| **k** | System Kill | Gracefully shuts down Controller, Broker, and all Workers |

## ✅ Testing

Comprehensive unit tests ensure correctness, race-safety, and image assembly validation.

```bash
# Test the distributed evolution logic
go test ./tests -v -run TestGol/-1\$

# Verify PGM output/assembly correctness
go test ./tests -v -run TestPgm/-1\$

# Check alive cell count reporting via RPC
go test ./tests -v -run TestAlive

# Run race condition detector (internal concurrency)
go test ./tests -v -race
```
