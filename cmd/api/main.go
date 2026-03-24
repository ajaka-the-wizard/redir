package main

import (
	"log"

	"github.com/ajaka-the-wizard/redir/internal"
)

func main() {
	if err := internal.Listen(); err != nil {
		log.Fatalf("couldn't bind to port; %v", err)
	}
}
