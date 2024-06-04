package main

import (
	"log"

	"github.com/xMajkel/x-kom-unboxer/pkg/roller"
	"github.com/xMajkel/x-kom-unboxer/pkg/utility/config"
)

func main() {
	err := config.ConfigInit()
	if err != nil {
		log.Printf("[ERROR] %+v\n", err)
		return
	}

	r := roller.Roller{}
	r.Start()
}
