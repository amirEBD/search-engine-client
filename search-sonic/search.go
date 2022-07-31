package main

import (
	"fmt"
	"github.com/expectedsh/go-sonic/sonic"
	"time"
)

const pswd = "password"

func main() {
	search, err := sonic.NewSearch("localhost", 1491, pswd)
	if err != nil {
		panic(err)
	}

	t0 := time.Now()
	results, _ := search.Query("messages", "default", "سلام", 10000, 0, sonic.LangAutoDetect)

	fmt.Println("time: ", time.Since(t0))
	fmt.Println("results count: ", len(results))
}
