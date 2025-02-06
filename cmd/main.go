//go:generate cp settings.yaml ./build/settings.yaml

package main

import (
	"fmt"
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

		addr := fmt.Sprintf(":%d", ctx.Settings.ListenPort)

		if err := r.RunTLS(addr, ctx.Settings.CertificateFile, ctx.Settings.PrivateKeyFile); err != nil {
			panic(err)
		}
	}()

	wg.Wait()

	ctx.Db.Close()
}
