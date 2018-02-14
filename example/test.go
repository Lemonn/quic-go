package main

import (
	"flag"
	"io"
	"net/http"
	"sync"

	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/lucas-clemente/quic-go/h2quic"
	"github.com/lucas-clemente/quic-go/internal/utils"
	"time"
)

func main() {
	verbose := flag.Bool("v", false, "verbose")
	flag.Parse()
	urls := flag.Args()

	if *verbose {
		utils.SetLogLevel(utils.LogLevelDebug)
	} else {
		utils.SetLogLevel(utils.LogLevelInfo)
	}
	utils.SetLogTimeFormat("")

	tr := &h2quic.RoundTripper{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	hclient := &http.Client{
		Transport: tr,
	}

	var wg sync.WaitGroup
	wg.Add(len(urls))
	addr := urls[0]
	utils.Infof("GET %s", addr)
	b := &bytes.Buffer{}
	quit := make(chan struct{})

	go func(addr string) {
		rsp, err := hclient.Get(addr)
		if err != nil {
			panic(err)
		}
		utils.Infof("Got response for %s: %#v", addr, rsp)
		_, err = io.Copy(b, rsp.Body)
		if err != nil {
			panic(err)
		}
		quit <- struct{}{}
		wg.Done()
	}(addr)

	go func() {
		time.Sleep(1000)
		for true {
			select {
			case <-quit:
				break
			default:
				b.Reset()
			}
		}

	}()
	wg.Wait()
}
