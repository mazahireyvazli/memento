package main

import (
	"time"

	"github.com/mazahireyvazli/memento"
)

func main() {
	memento.RunServer(&memento.MementoServerConfig{
		MementoConfig: &memento.MementoConfig{
			ShardNum:       1 << 10,
			ShardCapHint:   1 << 16,
			EntryExpiresIn: time.Minute * 10,
		},
		Port: 3000,
	})
}
