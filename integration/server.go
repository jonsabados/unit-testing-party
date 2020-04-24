package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NYTimes/gizmo/server/kit"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"google.golang.org/grpc"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Generation string

const (
	// see https://www.cnn.com/2013/11/06/us/baby-boomer-generation-fast-facts/index.html
	Greatest   Generation = "Greatest"     // 1924 and earlier
	Silent     Generation = "Silent"       // 1925 - 1945
	BabyBoomer Generation = "Baby Boomer"  // 1946 - 1964
	GenX       Generation = "Generation X" // 1965 - 1980
	Millennial Generation = "Millennial"   // 1981 - 1996
	GenZ       Generation = "Generation Z" // 1997 +
)

type RemoteEmployee struct {
	Status string `json:"status"`
	Data   struct {
		ID             string `json:"id"`
		EmployeeName   string `json:"employee_name"`
		EmployeeSalary string `json:"employee_salary"`
		EmployeeAge    string `json:"employee_age"`
		ProfileImage   string `json:"profile_image"`
	} `json:"data"`
}

type Employee struct {
	ID         string     `json:"id"`
	Name       string     `json:"employee_name"`
	Age        int        `json:"age"`
	Generation Generation `json:"generation"`
}

type Error struct {
	Message string `json:"message"`
}

type SomeServer struct {
	HttpClient *http.Client
}

func (s *SomeServer) EmployeeEndpoint(ctx context.Context, req interface{}) (interface{}, error) {
	employeeID := req.(int)

	url := fmt.Sprintf("http://dummy.restapiexample.com/api/v1/employee/%d", employeeID)
	_ = kit.LogDebugf(ctx, "fetching url %s", url)
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, kit.NewJSONStatusResponse(Error{err.Error()}, http.StatusInternalServerError)
	}
	res, err := s.HttpClient.Do(r)
	if err != nil {
		// this is going to be nearly impossible to test without black magic. What if we needed to do something
		// important in the error case?
		return nil, kit.NewJSONStatusResponse(Error{err.Error()}, http.StatusInternalServerError)
	}
	defer res.Body.Close()
	// 401 seems to be the response code from restapiexample for not found, wtf but OK
	if res.StatusCode == http.StatusUnauthorized {
		_, _ = ioutil.ReadAll(res.Body)
		return nil, kit.NewJSONStatusResponse(Error{"employee not found"}, http.StatusNotFound)
	} else if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		_ = kit.LogErrorf(ctx, "unexpected response code, got %d with body %s", res.StatusCode, string(body))
		return nil, kit.NewJSONStatusResponse(Error{"remote error encountered"}, http.StatusInternalServerError)
	}

	remote := new(RemoteEmployee)
	err = json.NewDecoder(res.Body).Decode(remote)
	if err != nil {
		return nil, kit.NewJSONStatusResponse(Error{err.Error()}, http.StatusInternalServerError)
	}

	age, err := strconv.Atoi(remote.Data.EmployeeAge)
	if err != nil {
		_ = kit.LogErrorf(ctx, "error parsing age, result: %+v", remote)
		return nil, kit.NewJSONStatusResponse(Error{err.Error()}, http.StatusInternalServerError)
	}

	ret := Employee{
		ID:   remote.Data.ID,
		Name: remote.Data.EmployeeName,
		Age:  age,
	}

	birthYear := time.Now().Year() - age

	// this is gonna -suck- test test, so it probably won't happen
	if birthYear <= 1924 {
		ret.Generation = Greatest
	} else if birthYear <= 1945 {
		ret.Generation = Silent
	} else if birthYear <= 1964 {
		ret.Generation = BabyBoomer
	} else if birthYear <= 1980 {
		ret.Generation = GenX
	} else if birthYear <= 1996 {
		ret.Generation = Millennial
	} else {
		ret.Generation = GenZ
	}

	return ret, nil
}

func getRequestID(_ context.Context, r *http.Request) (request interface{}, err error) {
	id, err := strconv.Atoi(kit.Vars(r)["id"])
	if err != nil {
		return nil, kit.NewJSONStatusResponse(Error{"employee not found"}, http.StatusNotFound)
	}

	return id, nil
}

func (s *SomeServer) Middleware(next endpoint.Endpoint) endpoint.Endpoint {
	return next
}

func (s *SomeServer) HTTPMiddleware(next http.Handler) http.Handler {
	return next
}

func (s *SomeServer) HTTPOptions() []kithttp.ServerOption {
	return []kithttp.ServerOption{}
}

func (s *SomeServer) HTTPRouterOptions() []kit.RouterOption {
	return []kit.RouterOption{}
}

func (s *SomeServer) HTTPEndpoints() map[string]map[string]kit.HTTPEndpoint {
	return map[string]map[string]kit.HTTPEndpoint{
		"/employee/{id}": {
			http.MethodGet: {
				Endpoint: s.EmployeeEndpoint,
				Decoder:  getRequestID,
			},
		},
	}
}

func (s *SomeServer) RPCMiddleware() grpc.UnaryServerInterceptor {
	return nil
}

func (s *SomeServer) RPCServiceDesc() *grpc.ServiceDesc {
	return nil
}

func (s *SomeServer) RPCOptions() []grpc.ServerOption {
	return nil
}
