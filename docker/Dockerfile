# ==================================================================================================================================
# This image specializes in providing a pipeline-as-a-container for Golang applications.
# It executes the following steps using a magefile (build.go):
#  - Vetting
#  - Security scan using Trivy
# Arguments:
#  - FILE_SRC: the name of the source file
# ==================================================================================================================================

FROM toremo/golang-build:latest AS build

ARG FILE_SRC
ARG FILE_DST

COPY . ./
RUN go build -o ${FILE_DST} ${FILE_SRC}

FROM alpine:latest

RUN adduser -D runner
RUN apk add --no-cache gcompat
COPY --from=build --chown=runner:runner /build/bin/app /usr/local/bin
USER runner 

ENTRYPOINT [ "${FILE_DST}" ]
