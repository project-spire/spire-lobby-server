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

	defer ctx.D.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		listenAddr := fmt.Sprintf(":%d", ctx.S.ListenPort)
		if err := r.RunTLS(listenAddr, ctx.S.CertificateFile, ctx.S.PrivateKeyFile); err != nil {
			panic(err)
		}
	}()

	wg.Wait()
}
