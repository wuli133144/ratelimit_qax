package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	//"sync/atomic"
	"time"
)

var (
	httpClient *http.Client
)

// init HTTPClient
func init() {
	httpClient = createHTTPClient()
}

const (
	MaxIdleConns        int = 4
	MaxIdleConnsPerHost int = 4
	IdleConnTimeout     int = 60
)

// createHTTPClient for connection re-use
func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 20 * time.Second,
			}).DialContext,
			MaxIdleConns:        MaxIdleConns,
			MaxIdleConnsPerHost: MaxIdleConnsPerHost,
			IdleConnTimeout:     time.Duration(IdleConnTimeout) * time.Second,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}
	return client
}

func init() {
	httpClient = createHTTPClient()
}


func main() {

	//s:=NewSafeTimerScheduel()
	//
	//go func() {
	//	for {
	//		df := <-s.GetTriggerChannel()
	//		df.Call()
	//		//atomic.AddInt64(&tt, -1)
	//	}
	//}()
	//
	//s.CreateTimer(100,test,[]interface{}{1,2})

	t := CreateRateLimit(2000, 100)
	for {
		pre := time.Now()
		r := t.TakeT()
		fmt.Println(r.Sub(pre))
	}
}
