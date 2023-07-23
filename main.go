package main

import (
	"fmt"
	"github.com/tad17/tad/fiber"
)

func main() {
	app := fiber.NewFiber()
	fmt.Printf("init fiber: %v\n", app)
}