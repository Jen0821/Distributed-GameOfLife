package gol

import (
	"fmt"
	"log"
	"net/rpc"
	"sync"
	"time"

	"uk.ac.bris.cs/gameoflife/server"
	"uk.ac.bris.cs/gameoflife/stubs"
	"uk.ac.bris.cs/gameoflife/util"
)

type distributorChannels struct {
	events     chan<- Event
	ioCommand  chan<- ioCommand
	ioIdle     <-chan bool
	ioFilename chan<- string
	ioOutput   chan<- uint8
	ioInput    <-chan uint8
	keyPresses <-chan rune
}

func copyWorld(w [][]byte) [][]byte {
	h := len(w)
	out := make([][]byte, h)
	for y := 0; y < h; y++ {
		out[y] = make([]byte, len(w[y]))
		copy(out[y], w[y])
	}
	return out
}

func savePGM(p Params, c distributorChannels, w [][]byte, turn int) {
	c.ioCommand <- ioOutput
	fn := fmt.Sprintf("%dx%dx%d", p.ImageWidth, p.ImageHeight, turn)
	c.ioFilename <- fn
	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			c.ioOutput <- w[y][x]
		}
	}
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle
	c.events <- ImageOutputComplete{CompletedTurns: turn, Filename: fn}
}

func distributor(p Params, c distributorChannels) {
	c.ioCommand <- ioInput
	c.ioFilename <- fmt.Sprintf("%dx%d", p.ImageWidth, p.ImageHeight)

	world := make([][]byte, p.ImageHeight)
	for y := 0; y < p.ImageHeight; y++ {
		row := make([]byte, p.ImageWidth)
		for x := 0; x < p.ImageWidth; x++ {
			row[x] = <-c.ioInput
		}
		world[y] = row
	}

	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			if world[y][x] == 0xFF {
				c.events <- CellFlipped{CompletedTurns: 0, Cell: util.Cell{X: x, Y: y}}
			}
		}
	}
	c.events <- StateChange{CompletedTurns: 0, NewState: Executing}

	var mu sync.Mutex
	turn := 0
	paused := false

	ticker := time.NewTicker(2 * time.Second)
	stop := make(chan struct{})
	go func() {
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				mu.Lock()
				cnt := 0
				for y := 0; y < p.ImageHeight; y++ {
					for x := 0; x < p.ImageWidth; x++ {
						if world[y][x] == 0xFF {
							cnt++
						}
					}
				}
				ct := turn
				mu.Unlock()
				c.events <- AliveCellsCount{CompletedTurns: ct, CellsCount: cnt}
			}
		}
	}()

	sliceH := p.ImageHeight / p.Threads
	results := make([]chan [][]byte, p.Threads)
	for i := 0; i < p.Threads; i++ {
		results[i] = make(chan [][]byte, 1)
	}

	for turn < p.Turns {
		select {
		case k := <-c.keyPresses:
			switch k {
			case 'p':
				paused = !paused
				mu.Lock()
				ct := turn
				mu.Unlock()
				if paused {
					c.events <- StateChange{CompletedTurns: ct, NewState: Paused}
				} else {
					c.events <- StateChange{CompletedTurns: ct, NewState: Executing}
				}
			case 's':
				mu.Lock()
				wc := copyWorld(world)
				ct := turn
				mu.Unlock()
				savePGM(p, c, wc, ct)
			case 'q':
				ticker.Stop()
				close(stop)
				mu.Lock()
				ct := turn
				al := make([]util.Cell, 0, 1024)
				for y := 0; y < p.ImageHeight; y++ {
					for x := 0; x < p.ImageWidth; x++ {
						if world[y][x] == 0xFF {
							al = append(al, util.Cell{X: x, Y: y})
						}
					}
				}
				wc := copyWorld(world)
				mu.Unlock()
				c.events <- FinalTurnComplete{CompletedTurns: ct, Alive: al}
				savePGM(p, c, wc, ct)
				c.events <- StateChange{CompletedTurns: ct, NewState: Quitting}
				close(c.events)
				return
			}
		default:
			if paused {
				time.Sleep(10 * time.Millisecond)
				continue
			}

			mu.Lock()
			cur := copyWorld(world)
			ct := turn
			mu.Unlock()

			results := make([][][]byte, p.Threads)
			if p.StartDistributed {
				var wg sync.WaitGroup
				for w := 0; w < p.Threads; w++ {
					wg.Add(1)
					go func(workerIdx int) {
						defer wg.Done()
						startY := workerIdx * sliceH
						h := sliceH
						if workerIdx == p.Threads-1 {
							h += p.ImageHeight % p.Threads
						}

						var haloUpper, haloLower []byte
						if workerIdx > 0 {
							lastWorkerEndY := (workerIdx-1)*sliceH + sliceH - 1
							if workerIdx-1 == p.Threads-1 {
								lastWorkerEndY = p.ImageHeight - 1
							}
							haloUpper = cur[lastWorkerEndY]
						}
						if workerIdx < p.Threads-1 {
							nextWorkerStartY := (workerIdx + 1) * sliceH
							haloLower = cur[nextWorkerStartY]
						}

						req := server.WorkerRequest{
							StartY:      startY,
							Height:      h,
							ImageWidth:  p.ImageWidth,
							ImageHeight: p.ImageHeight,
							World:       cur[startY : startY+h],
							Turn:        ct,
							HaloUpper:   haloUpper,
							HaloLower:   haloLower,
						}
						var res server.WorkerResponse

						client, err := rpc.Dial("tcp", p.WorkerAddrs[workerIdx])
						if err != nil {
							log.Fatalf("Worker %d connection failed: %v", workerIdx, err)
						}
						// Call RPC method (matches registered name in server.go)
						err = client.Call(stubs.DoWorkerHandler, req, &res)
						if err != nil {
							log.Fatalf("Worker %d RPC call failed: %v", workerIdx, err)
						}
						client.Close()

						results[workerIdx] = res.Result
					}(w)
				}
				wg.Wait()
			} else {
				var wg sync.WaitGroup
				for w := 0; w < p.Threads; w++ {
					wg.Add(1)
					go func(workerIdx int) {
						defer wg.Done()
						startY := workerIdx * sliceH
						h := sliceH
						if workerIdx == p.Threads-1 {
							h += p.ImageHeight % p.Threads
						}

						results[workerIdx] = Step(p, cur, startY, h)
					}(w)
				}
				wg.Wait()
			}

			next := make([][]byte, p.ImageHeight)
			for i := range next {
				next[i] = make([]byte, p.ImageWidth)
			}

			for w := 0; w < p.Threads; w++ {
				startY := w * sliceH
				h := sliceH
				if w == p.Threads-1 {
					h += p.ImageHeight % p.Threads
				}
				go func(wid, sy, hgt int) {
					res := Step(p, cur, sy, hgt)
					results[wid] = res
				}(w, startY, h)
			}

			for w := 0; w < p.Threads; w++ {
				startY := w * sliceH
				h := sliceH
				if w == p.Threads-1 {
					h += p.ImageHeight % p.Threads
				}
				part := results[w]
				for y := 0; y < h; y++ {
					copy(next[startY+y], part[y])
				}
			}

			flips := make([]util.Cell, 0, 2048)
			for y := 0; y < p.ImageHeight; y++ {
				for x := 0; x < p.ImageWidth; x++ {
					if cur[y][x] != next[y][x] {
						flips = append(flips, util.Cell{X: x, Y: y})
					}
				}
			}

			mu.Lock()
			world = next
			turn = ct + 1
			mu.Unlock()

			if len(flips) > 0 {
				c.events <- CellsFlipped{CompletedTurns: turn, Cells: flips}
			}
			c.events <- TurnComplete{CompletedTurns: turn}
		}
	}

	ticker.Stop()
	close(stop)

	mu.Lock()
	ct := turn
	al := make([]util.Cell, 0, 1024)
	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			if world[y][x] == 0xFF {
				al = append(al, util.Cell{X: x, Y: y})
			}
		}
	}
	wc := copyWorld(world)
	mu.Unlock()

	c.events <- FinalTurnComplete{CompletedTurns: ct, Alive: al}
	savePGM(p, c, wc, ct)
	c.events <- StateChange{CompletedTurns: ct, NewState: Quitting}
	close(c.events)
}
