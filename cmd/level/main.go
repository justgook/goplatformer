package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log/slog"
	"os"
)

func main() {
  	outPath := flag.String("o", "unknown.gob", "output location")
	flag.Parse()

	args := flag.Args()
        slog.Info(fmt.Sprintf("BUILDIN from: %v", args))
	file, _ := os.Create(*outPath)
	defer file.Close()

	output := "hello World"
	encoder := gob.NewEncoder(file)
	encoder.Encode(output)
}

