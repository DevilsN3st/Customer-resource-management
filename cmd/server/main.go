package main

import (
	"fmt"
	"github.com/icrxz/crm-api-core/internal"
	"time"
)

func main() {
	time.Local = time.UTC
	fmt.Println("Starting crm-core app!!")

	if err := internal.RunApp(); err != nil {
		panic(err)
	}
}
