package server

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

type WorkerRequest struct {
	StartY      int
	Height      int
	ImageWidth  int
	ImageHeight int
	World       [][]byte
	Turn        int
	HaloUpper   []byte
	HaloLower   []byte
}

type WorkerResponse struct {
	Result [][]byte
}

type GameOfLifeServer struct{}

func (g *GameOfLifeServer) DoWork(req WorkerRequest, res *WorkerResponse) error {
	height := req.Height
	width := req.ImageWidth
	world := req.World

	tempWorld := make([][]byte, height+2)
	if len(req.HaloUpper) > 0 {
		tempWorld = append(tempWorld, req.HaloUpper)
	}
	tempWorld = append(tempWorld, world...)
	if len(req.HaloLower) > 0 {
		tempWorld = append(tempWorld, req.HaloLower)
	}
	tempStartIdx := 0
	if len(req.HaloUpper) > 0 {
		tempStartIdx = 1
	}

	next := make([][]byte, height)
	for i := range next {
		next[i] = make([]byte, width)
	}

	for y := 0; y < height; y++ {
		tempY := tempStartIdx + y
		for x := 0; x < width; x++ {
			sum := 0
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if dx == 0 && dy == 0 {
						continue
					}
					nx := (x + dx + width) % width
					ny := tempY + dy
					if ny >= 0 && ny < len(tempWorld) {
						if tempWorld[ny][nx] == 0xFF {
							sum++
						}
					}
				}
			}
			cur := world[y][x]
			nextVal := byte(0x00)
			if cur == 0xFF {
				if sum == 2 || sum == 3 {
					nextVal = 0xFF
				}
			} else if sum == 3 {
				nextVal = 0xFF
			}
			next[y][x] = nextVal
		}
	}

	res.Result = next
	return nil
}

func RunServer() {
	server := new(GameOfLifeServer)
	err := rpc.Register(server)
	if err != nil {
		log.Fatalf("Error registering server: %v", err)
	}

	listener, err := net.Listen("tcp", ":8030")
	if err != nil {
		log.Fatalf("Listen error: %v", err)
	}

	fmt.Println("Server running on port 8030...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
