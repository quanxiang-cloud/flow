CONF ?=$(shell pwd)/configs
.PHONY: run-flow
generate:
	dapr run --app-id flow --app-port 9080 --components-path ${CONF}/samples -- go run cmd/main.go --config ${CONF}/config.yml