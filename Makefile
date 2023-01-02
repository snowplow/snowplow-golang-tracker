.PHONY: all format lint tidy test goveralls clean

# -----------------------------------------------------------------------------
#  CONSTANTS
# -----------------------------------------------------------------------------

build_dir = build

coverage_dir  = $(build_dir)/coverage
coverage_out  = $(coverage_dir)/coverage.out
coverage_html = $(coverage_dir)/coverage.html

# -----------------------------------------------------------------------------
#  BUILDING
# -----------------------------------------------------------------------------

all:
	GO111MODULE=on go build ./...

# -----------------------------------------------------------------------------
#  FORMATTING
# -----------------------------------------------------------------------------

format:
	GO111MODULE=on go fmt ./...
	GO111MODULE=on gofmt -s -w .

lint:
	GO111MODULE=on go install golang.org/x/lint/golint@latest
	GO111MODULE=on golint ./...

tidy:
	GO111MODULE=on go mod tidy

# -----------------------------------------------------------------------------
#  TESTING
# -----------------------------------------------------------------------------

test:
	mkdir -p $(coverage_dir)
	GO111MODULE=on go install golang.org/x/tools/cmd/cover@latest
	GO111MODULE=on go test ./... -tags test -v -covermode=count -coverprofile=$(coverage_out)
	GO111MODULE=on go tool cover -html=$(coverage_out) -o $(coverage_html)

goveralls: test
	GO111MODULE=on go install github.com/mattn/goveralls@latest
	GO111MODULE=on goveralls -coverprofile=$(coverage_out) -service=github

# -----------------------------------------------------------------------------
#  CLEANUP
# -----------------------------------------------------------------------------

clean:
	rm -rf $(build_dir)
