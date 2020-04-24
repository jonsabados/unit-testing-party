package unit

import (
	"context"
	"errors"
	"fmt"
	"github.com/NYTimes/gizmo/server/kit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEmployeeEndpoint_ErrorFetchingEmployee(t *testing.T) {
	asserter := assert.New(t)

	fetcher := &MockEmployeeFetcher{}

	fetcher.On("FetchEmployee", mock.Anything, 2).Return(nil, errors.New("testing FTW"))

	testInstance := SomeServer{
		EmployeeFetcher: fetcher,
		EmployeeMapper: func(employee *RemoteEmployee) (employee2 *Employee, err error) {
			asserter.Fail("we should not have reached this point")
			return nil, nil
		},
	}
	srv := kit.NewServer(&testInstance)
	ts := httptest.NewServer(srv)
	defer ts.Close()

	status, body := doRequest(ts.URL, "/employee/2")
	asserter.Equal(500, status)
	asserter.Equal("{\"message\":\"something terrible happened\"}", body)
}

func TestEmployeeEndpoint_ErrorMappingEmployee(t *testing.T) {
	asserter := assert.New(t)

	fetcher := &MockEmployeeFetcher{}

	expectedRemoteEmployee := &RemoteEmployee{
		Status: "whatever",
	}

	fetcher.On("FetchEmployee", mock.Anything, 2).Return(expectedRemoteEmployee, nil)

	testInstance := SomeServer{
		EmployeeFetcher: fetcher,
		EmployeeMapper: func(employee *RemoteEmployee) (*Employee, error) {
			asserter.Equal(expectedRemoteEmployee, employee)
			return nil, errors.New("KaBOOM")
		},
	}
	srv := kit.NewServer(&testInstance)
	ts := httptest.NewServer(srv)
	defer ts.Close()

	status, body := doRequest(ts.URL, "/employee/2")
	asserter.Equal(500, status)
	asserter.Equal("{\"message\":\"something terrible happened\"}", body)
}

func TestEmployeeEndpoint_NotFound(t *testing.T) {
	asserter := assert.New(t)

	fetcher := &MockEmployeeFetcher{}

	fetcher.On("FetchEmployee", mock.Anything, 2).Return(nil, nil)

	testInstance := SomeServer{
		EmployeeFetcher: fetcher,
		EmployeeMapper: func(employee *RemoteEmployee) (*Employee, error) {
			asserter.Fail("we should not have reached this point")
			return nil, nil
		},
	}
	srv := kit.NewServer(&testInstance)
	ts := httptest.NewServer(srv)
	defer ts.Close()

	status, body := doRequest(ts.URL, "/employee/2")
	asserter.Equal(404, status)
	asserter.Equal("{\"message\":\"employee not found\"}", body)
}

func TestEmployeeEndpoint_HappyPath(t *testing.T) {
	asserter := assert.New(t)

	fetcher := &MockEmployeeFetcher{}

	expectedRemoteEmployee := &RemoteEmployee{
		Status: "whatever",
	}

	result := Employee{
		ID:         "123",
		Name:       "Bob McTester",
		Age:        21,
		Generation: "DrinksRUs",
	}

	fetcher.On("FetchEmployee", mock.Anything, 2).Return(expectedRemoteEmployee, nil)

	testInstance := SomeServer{
		EmployeeFetcher: fetcher,
		EmployeeMapper: func(employee *RemoteEmployee) (*Employee, error) {
			asserter.Equal(expectedRemoteEmployee, employee)
			return &result, nil
		},
	}
	srv := kit.NewServer(&testInstance)
	ts := httptest.NewServer(srv)
	defer ts.Close()

	status, body := doRequest(ts.URL, "/employee/2")
	fetcher.AssertExpectations(t)
	asserter.Equal(200, status)
	asserter.Equal("{\"id\":\"123\",\"employee_name\":\"Bob McTester\",\"age\":21,\"generation\":\"DrinksRUs\"}\n", body)
}

func doRequest(apiBase string, path string) (int, string) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", apiBase, path), nil)
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

// yay, programmable mocks via testify!
type MockEmployeeFetcher struct {
	mock.Mock
}

func (m *MockEmployeeFetcher) FetchEmployee(ctx context.Context, employeeID int) (*RemoteEmployee, error) {
	args := m.Called(ctx, employeeID)
	ret := args.Get(0)
	if ret == nil {
		return nil, args.Error(1)
	}
	return ret.(*RemoteEmployee), args.Error(1)
}
