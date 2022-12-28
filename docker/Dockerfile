FROM ghcr.io/spolab/golang-runtime:v0.0.2

ARG APP
ARG DESC

LABEL org.opencontainers.image.source=https://github.com/spolab/petclinic-dapr
LABEL org.opencontainers.image.description=${DESC}
LABEL org.opencontainers.image.licenses=MIT

COPY ${APP} ./app

ENTRYPOINT [ "./actor" ]
