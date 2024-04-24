default: install

install: 
	go install
	
build: 
	go build

testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

doc:
	go generate ./...

fmt:
	go fmt ./...

deps:
	go mod tidy

.PHONY: build testacc doc fmt deps