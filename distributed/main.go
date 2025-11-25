package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"uk.ac.bris.cs/gameoflife/gol"
	"uk.ac.bris.cs/gameoflife/sdl"
	"uk.ac.bris.cs/gameoflife/server"
	"uk.ac.bris.cs/gameoflife/util"
)

// main is the function called when starting Game of Life with 'go run .'

func splitWorkerAddrs(s string) []string {
	if s == "" {
		return nil
	}
	var addrs []string

	for _, addr := range strings.Split(s, ",") {
		trimmed := strings.TrimSpace(addr)
		if trimmed != "" {
			addrs = append(addrs, trimmed)
		}
	}
	return addrs
}

func main() {
	runtime.LockOSThread()
	var params gol.Params

	startServer := flag.Bool("server", false, "start the distributed worker server(not main process)")

	flag.BoolVar(
		&params.StartDistributed,
		"dist",
		false,
		"Enable distributed mode (use remote workers)")

	var workerAddrsStr string
	flag.StringVar(
		&workerAddrsStr,
		"workers",
		"",
		"Remote worker addresses (comma-separated, e.g. \"ip:port,ip:port\")")

	flag.IntVar(
		&params.Threads,
		"t",
		8,
		"Specify the number of worker threads to use. Defaults to 8.")

	flag.IntVar(
		&params.ImageWidth,
		"w",
		512,
		"Specify the width of the image. Defaults to 512.")

	flag.IntVar(
		&params.ImageHeight,
		"h",
		512,
		"Specify the height of the image. Defaults to 512.")

	flag.IntVar(
		&params.Turns,
		"turns",
		10000000000,
		"Specify the number of turns to process. Defaults to 10000000000.")

	headless := flag.Bool(
		"headless",
		false,
		"Disable the SDL window for running in a headless environment.")

	flag.Parse()

	if *startServer {
		server.RunServer()
		return
	}

	params.WorkerAddrs = splitWorkerAddrs(workerAddrsStr)

	log.Printf("[Main] %-10v %v", "Threads", params.Threads)
	log.Printf("[Main] %-10v %v", "Width", params.ImageWidth)
	log.Printf("[Main] %-10v %v", "Height", params.ImageHeight)
	log.Printf("[Main] %-10v %v", "Turns", params.Turns)

	keyPresses := make(chan rune, 10)
	events := make(chan gol.Event, 1000)

	go sigint()

	go gol.Run(params, events, keyPresses)
	if !*headless {
		sdl.Run(params, events, keyPresses)
	} else {
		sdl.RunHeadless(events)
	}
}

func sigint() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
	var exit atomic.Bool
	for range sigint {
		if exit.Load() {
			log.Printf("[Main] %v Force quit by the user", util.Yellow("WARN"))
			os.Exit(0)
		} else {
			log.Printf("[Main] %v Press Ctrl+C again to force quit", util.Yellow("WARN"))
			exit.Store(true)
			go func() {
				time.Sleep(4 * time.Second)
				exit.Store(false)
			}()
		}
	}
}
