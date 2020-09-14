package unit

import (
	"github.com/jonsabados/unit-testing-party/unit/testutil"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRemoteEmployeeFetcher_HttpError(t *testing.T) {
	asserter := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
	}))
	ts.Close()

	testInstance := NewRemoteEmployeeFetcher(ts.URL)
	_, err := testInstance.FetchEmployee(testutil.NewTestContext(), 1)
	asserter.Error(err) // message will contain a random port so not gonna fuss with matching exact error
}

func TestRemoteEmployeeFetcher_NotFound(t *testing.T) {
	asserter := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		asserter.Equal("/api/v1/employee/1", r.RequestURI)
		bytes, err := ioutil.ReadFile("fixture/remote_employee_not_found.json")
		asserter.NoError(err)
		_, _ = w.Write(bytes)
	}))
	defer ts.Close()

	testInstance := NewRemoteEmployeeFetcher(ts.URL)
	res, err := testInstance.FetchEmployee(testutil.NewTestContext(), 1)
	asserter.NoError(err)
	asserter.Nil(res)
}

func TestRemoteEmployeeFetcher_Non200(t *testing.T) {
	asserter := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		asserter.Equal("/api/v1/employee/1", r.RequestURI)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("stuff went terribly wrong"))
	}))
	defer ts.Close()

	testInstance := NewRemoteEmployeeFetcher(ts.URL)
	res, err := testInstance.FetchEmployee(testutil.NewTestContext(), 1)
	asserter.EqualError(err, "unexpected response code, got 500 with body stuff went terribly wrong")
	asserter.Nil(res)
}

func TestRemoteEmployeeFetcher_GarbageRendered(t *testing.T) {
	asserter := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		asserter.Equal("/api/v1/employee/1", r.RequestURI)
		_, _ = w.Write([]byte("this isn't json, HA-HA!"))
	}))
	defer ts.Close()

	testInstance := NewRemoteEmployeeFetcher(ts.URL)
	res, err := testInstance.FetchEmployee(testutil.NewTestContext(), 1)
	asserter.Error(err)
	asserter.Nil(res)
}

func TestRemoteEmployeeFetcher_HappyPath(t *testing.T) {
	asserter := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		asserter.Equal("/api/v1/employee/1", r.RequestURI)
		bytes, err := ioutil.ReadFile("fixture/remote_employee.json")
		asserter.NoError(err)
		_, _ = w.Write(bytes)
	}))
	defer ts.Close()

	testInstance := NewRemoteEmployeeFetcher(ts.URL)
	res, err := testInstance.FetchEmployee(testutil.NewTestContext(), 1)
	asserter.NoError(err)
	asserter.Equal(&RemoteEmployee{
		Status: "success",
		Data: &struct {
			ID             int    `json:"id"`
			EmployeeName   string `json:"employee_name"`
			EmployeeSalary int    `json:"employee_salary"`
			EmployeeAge    int    `json:"employee_age"`
			ProfileImage   string `json:"profile_image"`
		}{
			ID:             1,
			EmployeeName:   "Tiger Nixon",
			EmployeeSalary: 320800,
			EmployeeAge:    61,
			ProfileImage:   "",
		},
	}, res)
}
