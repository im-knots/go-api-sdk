FROM golang:1.20.5 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN mkdir /build/
RUN go build -o /build/simple examples/simple/main.go
RUN go build -o /build/instrumented examples/instrumented/main.go


FROM gcr.io/distroless/base
COPY --from=builder /build/simple /bin/simple
COPY --from=builder /build/instrumented /bin/instrumented
