# Builder
FROM golang:1.20.12-alpine3.18 AS builderLinkShort
WORKDIR /cont
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN apk update && apk upgrade && \
    apk --update add git make
RUN go build -o linkshort ./cmd/main.go

FROM alpine:latest
RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    mkdir /app
WORKDIR /app

COPY --from=builderLinkShort ./cont/linkshort /app

CMD ["/app/linkshort"]