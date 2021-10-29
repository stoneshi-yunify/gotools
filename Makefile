GO_MODULE_ON=GO111MODULE=on
GO_ENV=${GO_MODULE_ON} GOOS=linux GOARCH=amd64 CGO_ENABLED=0

build-apiserver-sans-adder:
	${GO_ENV} go build -ldflags '-w -s' -v -tags netgo -o ./tmp/apiserver-sans-adder ./cmd/openssl/...