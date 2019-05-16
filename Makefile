.PHONY: all format lint test goveralls clean

# -----------------------------------------------------------------------------
#  CONSTANTS
# -----------------------------------------------------------------------------

src_dir       = tracker
build_dir     = build

coverage_dir  = $(build_dir)/coverage
coverage_out  = $(coverage_dir)/coverage.out
coverage_html = $(coverage_dir)/coverage.html

# -----------------------------------------------------------------------------
#  BUILDING
# -----------------------------------------------------------------------------

all: dep
	go build ./$(src_dir)

# -----------------------------------------------------------------------------
#  FORMATTING
# -----------------------------------------------------------------------------

format:
	go fmt ./$(src_dir)
	gofmt -s -w ./$(src_dir)

lint:
	go get -u github.com/golang/lint/golint
	golint ./$(src_dir)

# -----------------------------------------------------------------------------
#  TESTING
# -----------------------------------------------------------------------------

test:
	mkdir -p $(coverage_dir)
	go test ./$(src_dir) -tags test -v -covermode=count -coverprofile=$(coverage_out)
	go tool cover -html=$(coverage_out) -o $(coverage_html)

goveralls: test
	go get -u github.com/mattn/goveralls
	goveralls -coverprofile=$(coverage_out) -service=travis-ci

# -----------------------------------------------------------------------------
#  DEPENDENCIES
# -----------------------------------------------------------------------------

dep:
	dep ensure

# -----------------------------------------------------------------------------
#  CLEANUP
# -----------------------------------------------------------------------------

clean:
	rm -rf $(build_dir)
