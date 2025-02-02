//go:generate cp settings.yaml ./build/settings.yaml

package main

import (
	"sync"

	"spire/lobby/internal/core"
	"spire/lobby/internal/router"
)

func main() {
	ctx := core.NewContext()
	r := router.NewRouter(ctx)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := r.Run(":8080"); err != nil {
			panic(err)
		}
	}()

	wg.Wait()

	ctx.Db.Close()
}
