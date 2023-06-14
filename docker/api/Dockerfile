FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum main.go ./
COPY cmd/ ./cmd
COPY internal/ ./internal
RUN /bin/bash -c "go mod download && go build -o ./app && find . -not -name 'app' -delete"