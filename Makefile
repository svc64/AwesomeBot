.PHONY: all test clean

GOBIN = build/bin

all:
	./fakespace.sh go get github.com/golang/dep/cmd/dep
	./fakespace.sh $(GOBIN)/dep ensure
	./fakespace.sh ./build.sh

test: all
	./fakespace.sh go test -v ./...

vet: all
	./fakespace.sh go vet -shadow=false ./...

clean:
	rm -rf out
