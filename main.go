package main

import (
	"backend/app"
	"backend/config"
	"fmt"
)

func main() {
	fmt.Println("INICIO")
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("error config: %+v", err)
		return
	}
	app.Run(cfg)
}
