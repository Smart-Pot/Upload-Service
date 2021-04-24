package main

import (
	"log"
	"os"
	"os/signal"

	"uploadservice/cmd"

	"github.com/Smart-Pot/pkg"
)

func main() {
	if err := pkg.Config.ReadConfig(); err != nil {
		log.Fatal(err)
	}
	log.Println("Configurations are set")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		if err := cmd.Execute(); err != nil {
			log.Fatal(err)
		}
	}()
	sig := <-c
	log.Println("GOT SIGNAL: " + sig.String())
}
