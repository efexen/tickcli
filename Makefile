NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)

all: deps
	@mkdir -p bin/
	@echo "$(OK_COLOR)==> Building$(NO_COLOR)"
	@go build -o bin/tickcli

deps:
	@go get -d -v ./...
	@echo "$(OK_COLOR)==> Installing dependencies$(NO_COLOR)"
	@echo $(DEPS) | xargs -n1 go get -d
