package main

import (
	"github.com/NYTimes/gizmo/server/kit"
	"github.com/jonsabados/unit-testing-party/unit"
	"net/http"
)

func main() {
	svc := unit.SomeServer{
		EmployeeFetcher: unit.NewRemoteEmployeeFetcher("http://dummy.restapiexample.com"),
		EmployeeMapper: unit.NewEmployeeFactory(unit.MapBirthYear),
	}
	svr := kit.NewServer(&svc)

	panic(http.ListenAndServe("0:8080", svr))
}
