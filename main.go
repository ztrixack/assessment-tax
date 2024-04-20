package main

import (
	"github.com/ztrixack/assessment-tax/internal/api"
)

func main() {
	server := api.New()
	server.Listen()
}
