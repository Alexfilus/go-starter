default: run

 
test: test test-coverage test-integration

run:
	go run cmd/zample/main.go --configFile="cmd/zample/config.toml"


# It adds any missing module requirements necessary to build the current module’s packages and dependencies, and it
# removes requirements on modules that don’t provide any relevant packages. It also adds any missing entries to go.sum
# and removes unnecessary entries.
.PHONY: tidy
tidy:
	go mod tidy

##


# compiles existing programs, and creates its binaries into bin folder
.PHONY: build
build:
	go build -o bin/zample ./cmd/zample/main.go 

.PHONY: test
test:
	go test -v  ./...

# runs coverage tests and generates the coverage report
test-coverage:
	go test ./... -v -coverpkg=./...

# runs integration tests
test-integration:
	go test ./... -tags=integration ./...