FROM golang:1.23-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .
RUN go build -o ./bin/app ./cmd/main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/app /
COPY config/dev.yaml /dev.yaml
ENV CONFIG_PATH=/dev.yaml
EXPOSE 8085

CMD ["/app"]