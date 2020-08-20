TEST?=$$(go list ./... | grep -v '/vendor/')
VETARGS?=-all
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

default: test vet

tools:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.25.0

clean:
	rm -Rf $(CURDIR)/bin/*

build: clean vet
	go build -o $(CURDIR)/bin/terraform-provider-mackerel $(CURDIR)//main.go

test: vet
	TF_ACC= go test $(TEST) $(TESTARGS) -timeout=30s -parallel=4

testacc: vet
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

vet: fmt
	@echo "go vet $(VETARGS) ."
	@go vet $(VETARGS) . ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

lint:
	@echo "==> Checking source code against linters..."
	@golangci-lint run

.PHONY: default test vet fmt tools lint
