FROM alpine:3.16

ARG DIR_HOME=/build

RUN apk update && \
    apk upgrade && \
    apk add --no-cache go
RUN adduser -D -h ${DIR_HOME} builder
USER builder
WORKDIR ${DIR_HOME}
RUN ls -la
RUN go install -buildvcs=true github.com/magefile/mage@latest