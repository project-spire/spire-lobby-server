package main

import (
	"fmt"
	"spire/lobby/internal/route"
	"sync"

	"spire/lobby/internal/core"
)

func main() {
	ctx := core.NewContext()
	r := route.NewRouter(ctx)

	defer ctx.Close()

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
