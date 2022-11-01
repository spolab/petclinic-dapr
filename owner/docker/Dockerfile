FROM golang:1.19.2 AS build

RUN mkdir /build
WORKDIR /build
COPY * .
RUN go vet && \
    go build -o ./app ./...

FROM alpine:3.16

WORKDIR /
RUN apk update && \
    apk add --no-cache gcompat
COPY --from=build /build/app .
EXPOSE 3000
ENTRYPOINT [ "/app" ]