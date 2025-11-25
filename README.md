# 🌐 Distributed-GameOfLife-Go: Scalable GOL Engine

This repository hosts a high-performance, **distributed implementation** of **Conway's Game of Life (GoL)**. Built entirely in **Go (Golang)**, the system leverages an asynchronous, network-based architecture using **Remote Procedure Calls (RPC)** to distribute computation across multiple machines (simulated as AWS nodes), achieving massive scalability for large grid simulations.

The core design separates I/O and visualization (Local Controller) from the primary computational load (AWS Nodes/Workers).

---

## 📝 Final Coursework Report

The full detailed analysis of the parallel implementation, including performance benchmarks and design rationale, is available in the final report.

[![Report Cover Image: Click to Download PDF](./report.jpg)](./report.pdf)
### Direct Link
[Download the Full PDF Report Here](./report.pdf)

---

## ✨ Distributed System Architecture

The project is structured around a three-tier model to handle tasks, communication, and computation:

1.  **Local Controller (Client):** Handles user I/O (PGM load/save, SDL visualization, Keypresses). It sends commands to the **Broker** via RPC.
2.  **Broker / Distributor (Central Server):** The orchestrator. It receives commands, divides the game board into slices, and delegates tasks to the **Workers**. It does **not** contain any Game of Life evolution logic.
3.  **GOL Worker (Compute Node):** Runs on separate machines (AWS nodes). It receives a board slice, calculates the next state, and returns the result. It employs **Halo Exchange** (communicating border rows/columns with neighbors) to maintain calculation correctness efficiently.

---

## ⚙️ Running the Distributed System

This project requires simultaneous execution of the Broker, one or more Workers, and the Local Controller.

### Prerequisites

* Go (Golang) installed.
* SDL2 development libraries installed (for visualization on the Controller).

### Execution Sequence

The components must be started in the following order:

#### 1. Start the Broker / Distributor (Central Server)

This component listens for RPC requests from both the Controller and the Workers.

```bash
# Start the Broker on a specified port (e.g., 8080)
go run ./broker -listen :8080

#### 2. Start the GOL Worker Nodes

These instances connect to the Broker and wait for work slices. Start as many instances as desired for parallel computation.

# Start Worker 1 (connected to Broker at :8080)
go run ./worker -broker-addr :8080 -worker-id 1

# Start Worker 2 (on a different machine/terminal)
go run ./worker -broker-addr :8080 -worker-id 2

# ... and so on for N workers

#### 3. Start the Local Controller (Client/I/O)

This controls the entire simulation. Use the -headless flag to disable SDL visualization for performance testing.

# Start the Controller, connecting to the Broker and initiating the simulation
go run ./controller -broker-addr :8080 -headless -t 4
