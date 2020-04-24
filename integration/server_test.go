package integration

import (
	"fmt"
	"github.com/NYTimes/gizmo/server/kit"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEmployeeEndpoint(t *testing.T) {
	// we could fire up the server every test but the annoying first request fails always thing makes doing a single
	// instance way more practical
	cookieJar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		panic(err)
	}

	svc := SomeServer{
		HttpClient: &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           cookieJar, // the sample api were using always errors on the first request....
			Timeout:       0,
		},
	}
	svr := kit.NewServer(&svc)
	ts := httptest.NewServer(svr)
	defer ts.Close()

	doRequest := func(path string) (int, string) {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", ts.URL, path), nil)
		if err != nil {
			panic(err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}

		defer res.Body.Close()
		bytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		return res.StatusCode, string(bytes)
	}

	doRequest("/employee/1")

	testCases := []struct {
		desc         string
		toFetch      string
		wantedStatus int
		wantedBody   string
	}{
		{
			"happy path baby boomer",
			"/employee/1",
			200,
			`{"id":"1","employee_name":"Tiger Nixon","age":61,"generation":"Baby Boomer"}`,
		},
		{
			"happy path baby millennial",
			"/employee/5",
			200,
			`{"id":"5","employee_name":"Airi Satou","age":33,"generation":"Millennial"}`,
		},
		{
			"non numeric",
			"/employee/BLAH",
			404,
			`{"message":"employee not found"}`,
		},
		{
			"not found",
			"/employee/9999",
			404,
			`{"message":"employee not found"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			asserter := assert.New(t)

			status, content := doRequest(tc.toFetch)
			asserter.Equal(tc.wantedStatus, status)
			asserter.Equal(tc.wantedBody, strings.Trim(content, "\n"))
		})
	}
}
