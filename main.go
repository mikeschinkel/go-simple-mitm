package main

import (
	"fmt"
	"log"
	"simple-mitm/src"
	"sync"
)

func main() {
	simple_mitm.Initialize()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		p := simple_mitm.NewProxy()
		fmt.Printf("Starting proxy at port %s\n", simple_mitm.ProxyPort)
		if err := p.ListenAndServeTLS(simple_mitm.PemPath, simple_mitm.KeyPath); err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()

	go func() {
		api := simple_mitm.NewAPI()
		fmt.Printf("Starting API server at port %s\n", simple_mitm.APIServerPort)
		if err := api.ListenAndServeTLS(simple_mitm.PemPath, simple_mitm.KeyPath); err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()
	wg.Wait()
}
