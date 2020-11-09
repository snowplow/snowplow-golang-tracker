.PHONY: all stress format lint tidy test goveralls protoc-gen-go collector-protocol-gen clean

# -----------------------------------------------------------------------------
#  CONSTANTS
# -----------------------------------------------------------------------------

src_dir   = tracker
build_dir = build
cmd_dir   = cmd

coverage_dir  = $(build_dir)/coverage
coverage_out  = $(coverage_dir)/coverage.out
coverage_html = $(coverage_dir)/coverage.html

collector_protocol = grpc/collector.proto

stress_src_dir = $(cmd_dir)/stress

# -----------------------------------------------------------------------------
#  BUILDING
# -----------------------------------------------------------------------------

all:
	GO111MODULE=on go build ./$(src_dir)/...

stress:
	mkdir -p $(build_dir)
	GO111MODULE=on go build -o ./$(build_dir) ./$(stress_src_dir)

# -----------------------------------------------------------------------------
#  FORMATTING
# -----------------------------------------------------------------------------

format:
	GO111MODULE=on go fmt ./$(src_dir)/...
	GO111MODULE=on gofmt -s -w ./$(src_dir)
	GO111MODULE=on go fmt ./$(cmd_dir)/...
	GO111MODULE=on gofmt -s -w ./$(cmd_dir)

lint:
	GO111MODULE=on go get -u golang.org/x/lint/golint
	GO111MODULE=on golint ./$(src_dir)/...
	GO111MODULE=on golint ./$(cmd_dir)/...

tidy:
	GO111MODULE=on go mod tidy

# -----------------------------------------------------------------------------
#  TESTING
# -----------------------------------------------------------------------------

test:
	mkdir -p $(coverage_dir)
	GO111MODULE=on go test ./$(src_dir)/... -tags test -v -covermode=count -coverprofile=$(coverage_out)
	GO111MODULE=on go tool cover -html=$(coverage_out) -o $(coverage_html)

goveralls: test
	GO111MODULE=on go get -u github.com/mattn/goveralls
	GO111MODULE=on goveralls -coverprofile=$(coverage_out) -service=travis-ci

# -----------------------------------------------------------------------------
#  COLLECTOR GRPC PROTOCOLS
# -----------------------------------------------------------------------------

protoc-gen-go:
	GO111MODULE=on go get google.golang.org/protobuf/cmd/protoc-gen-go google.golang.org/grpc/cmd/protoc-gen-go-grpc

# Protoc Installation: https://grpc.io/docs/protoc-installation/
collector-protocol-gen: protoc-gen-go
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
	$(collector_protocol)

# -----------------------------------------------------------------------------
#  CLEANUP
# -----------------------------------------------------------------------------

clean:
	rm -rf $(build_dir)
