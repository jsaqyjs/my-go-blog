package main

import (
	"log"
	"my-go-blog/routers"
)

func main() {
	r := routers.InitRouter()
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
