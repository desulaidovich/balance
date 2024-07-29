package main

import "github.com/desulaidovich/balance/internal/balance"

func main() {
	if err := balance.Run(); err != nil {
		panic(err)
	}
}
