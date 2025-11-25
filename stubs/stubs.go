package stubs

var DoWorkerHandler = "GameOfLifeService.DoWorker"

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
