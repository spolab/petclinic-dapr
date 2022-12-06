FROM alpine:3.16 AS build

ARG DIR_BUILD=/build
RUN apk add --no-cache go
RUN mkdir -p ${DIR_BUILD}
COPY . .${DIR_BUILD}
WORKDIR /${DIR_BUILD}
RUN go vet ./...
RUN go build -o bin/owner owner/cmd/app.go

FROM alpine:3.16

ARG DIR_BUILD=/build
ARG DIR_BIN=/usr/bin
RUN apk add --no-cache gcompat
COPY --from=build ${DIR_BUILD}/bin/owner ${DIR_BIN}
ENTRYPOINT ["owner"]