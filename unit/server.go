package unit

import (
	"context"
	"github.com/NYTimes/gizmo/server/kit"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"google.golang.org/grpc"
	"net/http"
	"strconv"
)

type Error struct {
	Message string `json:"message"`
}

type SomeServer struct {
	EmployeeFetcher RemoteEmployeeFetcher
	EmployeeMapper  EmployeeConverter
}

func (s *SomeServer) EmployeeEndpoint(ctx context.Context, req interface{}) (interface{}, error) {
	employeeID := req.(int)

	remote, err := s.EmployeeFetcher.FetchEmployee(ctx, employeeID)
	if err != nil {
		_ = kit.LogErrorf(ctx, "error reading employee %+v", err)
		return nil, kit.NewJSONStatusResponse(Error{"something terrible happened"}, http.StatusInternalServerError)
	}
	if remote == nil {
		return nil, kit.NewJSONStatusResponse(Error{"employee not found"}, http.StatusNotFound)
	}

	ret, err := s.EmployeeMapper(remote)
	if err != nil {
		_ = kit.LogErrorf(ctx, "error mapping employee, result: %+v, err: %s", remote, err)
		return nil, kit.NewJSONStatusResponse(Error{"something terrible happened"}, http.StatusInternalServerError)
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
