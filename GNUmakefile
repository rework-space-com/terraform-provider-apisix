default: testacc

# Run acceptance tests
.PHONY: testacc

testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

build: 
	go install

doc:
	go get github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
	go generate ./...

fmt:
	go fmt ./...

deps:
	go mod tidy