FROM golang:1.16-alpine AS build_deps

RUN apk add --no-cache git
WORKDIR /workspace
COPY go.mod .
COPY go.sum .

RUN go mod download

FROM build_deps AS build
COPY . .
RUN CGO_ENABLED=0 go build -o speedtest -ldflags '-w -extldflags "-static"' *.go

FROM alpine:3.14.3
RUN apk add --no-cache ca-certificates
COPY --from=build /workspace/speedtest /usr/local/bin/speedtest

ENTRYPOINT ["speedtest"]
