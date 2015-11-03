package main

import (
	"fmt"

	"github.com/jbowens/muni-display/server/core"
	"github.com/octavore/naga/service"
)

func init() {
	service.BootPrintln = func(v ...interface{}) {
		fmt.Println(v...)
	}
}

func main() {
	var server core.Module
	service.Run(&server)
}
