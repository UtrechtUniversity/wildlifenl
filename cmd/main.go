package main

import (
	"log"

	"gitlab.com/bazzz/config"
	"github.com/UtrechtUniversity/wildlifenl"
)

func main() {
	cfg := new(wildlifenl.Configuration)
	config.LoadFromApplicationPath(cfg)
	log.Fatal(wildlifenl.Start(cfg))
}
