package gol

type Params struct {
	Turns            int
	Threads          int
	ImageWidth       int
	ImageHeight      int
	StartDistributed bool
	WorkerAddrs      []string
}

func Run(p Params, events chan<- Event, keyPresses <-chan rune) {
	ioCommand := make(chan ioCommand)
	ioIdle := make(chan bool)
	ioFilename := make(chan string)
	ioOutput := make(chan uint8)
	ioInput := make(chan uint8)

	ioCh := ioChannels{
		command:  ioCommand,
		idle:     ioIdle,
		filename: ioFilename,
		output:   ioOutput,
		input:    ioInput,
	}
	go startIo(p, ioCh)

	c := distributorChannels{
		events:     events,
		ioCommand:  ioCommand,
		ioIdle:     ioIdle,
		ioFilename: ioFilename,
		ioOutput:   ioOutput,
		ioInput:    ioInput,
		keyPresses: keyPresses,
	}
	distributor(p, c)
}

func Step(p Params, world [][]byte, startY, height int) [][]byte {
	next := make([][]byte, height)
	for i := range next {
		next[i] = make([]byte, p.ImageWidth)
	}
	for y := 0; y < height; y++ {
		gy := startY + y
		for x := 0; x < p.ImageWidth; x++ {
			s := 0
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if dx == 0 && dy == 0 {
						continue
					}
					nx := (x + dx + p.ImageWidth) % p.ImageWidth
					ny := (gy + dy + p.ImageHeight) % p.ImageHeight
					if world[ny][nx] == 0xFF {
						s++
					}
				}
			}
			if world[gy][x] == 0xFF {
				if s == 2 || s == 3 {
					next[y][x] = 0xFF
				} else {
					next[y][x] = 0x00
				}
			} else {
				if s == 3 {
					next[y][x] = 0xFF
				} else {
					next[y][x] = 0x00
				}
			}
		}
	}
	return next
}
