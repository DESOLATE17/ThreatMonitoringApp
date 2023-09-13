package main

import (
	"threat-monitoring/internal/api/handler"
)

func main() {
	r := handler.StartServer()
	r.Run()
}
