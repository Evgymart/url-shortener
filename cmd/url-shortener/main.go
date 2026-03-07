package main

import (
	"fmt"
	"urlshort/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
}
