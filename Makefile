.PHONY: genall ogen sqlc install-tools testunit

ogen:
	ogen --target ./internal/transport/http/httpgen --package httpgen --clean ./api/v1/openapi.yaml

sqlc:
	sqlc generate

mocks:
	mockery --log-level=debug

genall: ogen sqlc mocks

install-tools:
	go install -v github.com/ogen-go/ogen/cmd/ogen@latest
	go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
	go install github.com/vektra/mockery/v3@v3.6.1

testunit:
	go test ./internal/usecase