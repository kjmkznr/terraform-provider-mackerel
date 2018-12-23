TEST?=$$(go list ./... | grep -v '/vendor/')
VETARGS?=-all
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

default: test vet

tools:
	go get -u github.com/kardianos/govendor
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install

clean:
	rm -Rf $(CURDIR)/bin/*

build: clean vet
	govendor build -o $(CURDIR)/bin/terraform-provider-mackerel $(CURDIR)/builtin/bins/provider-mackerel/main.go

test: vet
	TF_ACC= go test $(TEST) $(TESTARGS) -timeout=30s -parallel=4

testacc: vet
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

vet: fmt
	@echo "go tool vet $(VETARGS) ."
	@go tool vet $(VETARGS) $$(ls -d */ | grep -v vendor) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

lint:
	@echo "==> Checking source code against linters..."
	@gometalinter ./

.PHONY: default test vet fmt tools lint
