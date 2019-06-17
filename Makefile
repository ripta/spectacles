
all: build

build:
	go build -v ./cmd/spectacles

build-race:
	go build -race -v ./cmd/spectacles

fmt:
	goimports -local k8s.io,github.com/ripta/spectacles -w .

run: build
	./spectacles --kubeconfig $$KUBECONFIG

test:
	go test ./...
