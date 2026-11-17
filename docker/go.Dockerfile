# Dev/test toolchain image for Go tasks (not for production).
FROM golang:1.26
ENV GOBIN=/usr/local/bin
RUN go install github.com/pressly/goose/v3/cmd/goose@latest \
 && go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest \
 && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b /usr/local/bin
WORKDIR /workspace
