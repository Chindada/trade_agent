package main

import (
	"fmt"

	"gitlab.tocraw.com/root/toc_trader/pkg/config"
)

func main() {
	a, _ := config.Get()
	fmt.Println(a.GetDBConfig())
	fmt.Println(a.GetServerConfig())
}
