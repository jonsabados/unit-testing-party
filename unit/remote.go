package unit

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NYTimes/gizmo/server/kit"
	"github.com/pkg/errors"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
)

// Note, I would prefer just to do a type that is a function for this since were really just passing behavior around,
// but structs containing dependencies is way more familiar for OO folks and this will give us a good thing to
// demonstrate mocking interfaces with testify
type RemoteEmployeeFetcher interface {
	FetchEmployee(ctx context.Context, employeeID int) (*RemoteEmployee, error)
}

type restEmployeeFetcher struct {
	apiURL string
	client *http.Client
}

func (r *restEmployeeFetcher) FetchEmployee(ctx context.Context, employeeID int) (*RemoteEmployee, error) {
	url := fmt.Sprintf("%s/api/v1/employee/%d", r.apiURL, employeeID)
	_ = kit.LogDebugf(ctx, "fetching url %s", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	res, err := r.client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		return nil, errors.New(fmt.Sprintf("unexpected response code, got %d with body %s", res.StatusCode, string(body)))
	}

	remote := new(RemoteEmployee)
	err = json.NewDecoder(res.Body).Decode(remote)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if remote.Data == nil {
		return nil, nil
	}
	return remote, nil
}

func NewRemoteEmployeeFetcher(apiURL string) RemoteEmployeeFetcher {
	cookieJar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		panic(err)
	}

	ret := &restEmployeeFetcher{
		apiURL: apiURL,
		client: &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           cookieJar,
			Timeout:       0,
		},
	}

	return ret
}
