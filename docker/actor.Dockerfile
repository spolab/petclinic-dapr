FROM ghcr.io/spolab/golang-runtime:v0.0.2

ARG NAME

COPY bin/${NAME}/actor ./actor

ENTRYPOINT [ "./actor" ]
