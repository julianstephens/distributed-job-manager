package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/julianstephens/distributed-task-scheduler/pkg"
	"github.com/julianstephens/distributed-task-scheduler/seeds"
)

func main() {
	godotenv.Load()
	handleArgs()
}

func handleArgs() {
	flag.Parse()
	args := flag.Args()

	if len(args) >= 1 {
		switch args[0] {
		case "seed":
			db, err := pkg.GetDB()
			if err != nil {
				log.Fatalf("unable to init dynamodb client, %v", err)
			}
			masterSeedCount := 10
			seeds.Execute(db, masterSeedCount, args[1:]...)
			os.Exit(0)
		}
	}
}
