package main

import (
	"log"

	"github.com/ajaka-the-wizard/redir/internal"
)

func main() {
	err := internal.Listen()
	if err != nil {
		log.Println("Couldnt bind to port")
		log.Println(err.Error())
	}
}
