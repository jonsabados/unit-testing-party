# Unit testing party

This is two different examples of a fairly simple app. One is written in a testable way with unit tests (in the unit package), and the other written in a way that is only really testable via integration tests (in the integration package).

Note, the unit example is not meant to say that integration tests are not needed, there should absolutely be some sort of integration test on things. But, it there are a whole lotta code paths that are way, way, way easier to test with unit tests and it is meant to demonstrate that.

## Running things
Execute tests: `go test ./...`

Start integration test only version server `go run integration/cmd/main.go`

Start version of the server that is unit testable `go run integration/cmd/main.go`