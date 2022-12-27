FROM ghcr.io/spolab/golang-runtime:v0.0.2

ARG NAME

LABEL org.opencontainers.image.source=https://github.com/spolab/petstore-dapr
LABEL org.opencontainers.image.description="dapr petstore ${NAME} actor"
LABEL org.opencontainers.image.licenses=MIT

COPY bin/${NAME}/actor ./actor

ENTRYPOINT [ "./actor" ]
