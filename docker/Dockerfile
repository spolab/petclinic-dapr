# ==================================================================================================================================
# This image specializes in providing a pipeline-as-a-container for Golang applications.
# It executes the following steps using a magefile (build.go):
#  - Vetting
#  - Security scan using Trivy
# Arguments:
#  - FILE_SRC: the name of the source file
# ==================================================================================================================================

FROM alpine:3.16 AS build

ARG FILE_SRC

RUN mkdir -p /build && \
    apk add --no-cache go
WORKDIR /build
COPY . ./
RUN go build -o bin/app ${FILE_SRC}

FROM alpine:3.16

RUN adduser -D runner
RUN apk add --no-cache gcompat
COPY --from=build --chown=runner:runner /build/bin/app /usr/local/bin
USER runner 
ENTRYPOINT ["/usr/local/bin/app"]