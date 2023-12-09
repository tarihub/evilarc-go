export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
LDFLAGS := -s -w

all: fmt build

build: evilarc

fmt:
	go fmt ./...

test:
	go test

vet:
	go vet ./...

evilarc:
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/evilarc .

clean:
	rm -f ./bin/evilarc
	rm -f *.zip
	rm -f *.tar.gz
	rm -f *.tar