FROM golang:1.19.3 AS build

ARG DIR_BUILD=/build
ARG FILE_OUT=bin/app

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

RUN mkdir -p ${DIR_BUILD}
WORKDIR ${DIR_BUILD}
COPY . ${DIR_BUILD}/
RUN cd cmd && go build -buildvcs=false -o ../bin/app .

# -----------------------------------------------------------------------------

FROM alpine:3.16

ARG DIR_BUILD=/build
ARG FILE_OUT=bin/app

RUN mkdir -p /app
WORKDIR /app
RUN apk update && \
    apk add --no-cache gcompat
COPY --from=build ${DIR_BUILD}/${FILE_OUT} .
EXPOSE 3000
ENTRYPOINT [ "/app/app" ]