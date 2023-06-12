FROM golang:1.20.5 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN mkdir /build/
RUN go build -o /build/simple examples/simple/main.go
RUN go build -o /build/instrumented examples/instrumented/main.go


FROM gcr.io/distroless/static-debian11
COPY --from=build /build/simple /simple
COPY --from=build /build/instrumented /simple
