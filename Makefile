GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "./vendor/*" -not -path "./.git/*")

test:
	go test -v -cover ./...
vet:
	go vet -v  -v ./...

mod-update:
	go get -v -u ./...
	go mod tidy
	go mod vendor

# lint runs the staticcheck and golint static analysis tools on all packages in the project.
lint:
	go install mvdan.cc/gofumpt@latest
	gofumpt -l -w ${GOFILES_NOVENDOR}
	# $(call check_command_exists,staticcheck) || go install honnef.co/go/tools/cmd/staticcheck@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...
	# $(call check_command_exists,golint) || go install -v golang.org/x/lint/golint@latest
	go install -v golang.org/x/lint/golint@latest
	golint ${GOFILES_NOVENDOR}

# check_command_exists is a helper function that checks if a command exists.
check_command_exists = $(shell command -v $(1) > /dev/null && echo "true" || echo "false")

ifeq ($(call check_command_exists,$(1)),false)
  $(error "$(1) command not found")
endif

# help prints a list of available targets and their descriptions.
help:
	@echo "Available targets:"
	@echo
	@echo "vet\t\t\t\tRun the Go vet static analysis tool on all packages in the project."
	@echo "lint\t\t\t\tRun the staticcheck and golint static analysis tools on all packages in the project."
	@echo "test\t\t\t\tRun all tests in the project."
	@echo "mod-update\t\t\tUpdate all dependencies in the project."
	@echo "help\t\t\t\tPrint this help message."
	@echo
	@echo "For more information, see the project README."
