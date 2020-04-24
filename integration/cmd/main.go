package main

import (
	"github.com/NYTimes/gizmo/server/kit"
	"github.com/jonsabados/unit-testing-party/integration"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
)

func main() {
	cookieJar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		panic(err)
	}

	svc := integration.SomeServer{
		HttpClient: &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           cookieJar, // the sample api were using always errors on the first request....
			Timeout:       0,
		},
	}
	svr := kit.NewServer(&svc)

	panic(http.ListenAndServe("0:8080", svr))
}
