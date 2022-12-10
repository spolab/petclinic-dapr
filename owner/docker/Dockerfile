# ==================================================================================================================================
# This image specializes in providing a pipeline-as-a-container for Golang applications.
# It executes the following steps using a magefile (build.go):
#  - Vetting
#  - Security scan using Trivy
# Arguments:
#  - FILE_MAIN: the file identifying the main programme
#  - FILE
# ==================================================================================================================================

ARG FILE_MAIN

FROM spolab/toolbox:latest AS build

COPY --chown=builder:builder . ./
RUN go build -o bin/actor owner/cmd/actor.go

FROM alpine:3.16

RUN adduser -D runner
RUN apk add --no-cache gcompat
COPY --from=build /build/bin/actor /usr/local/bin
USER runner 
ENTRYPOINT ["actor"]