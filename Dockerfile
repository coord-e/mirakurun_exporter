# syntax=docker.io/docker/dockerfile:1

FROM gcr.io/distroless/static-debian11:nonroot

ARG BIN_DIR
ARG TARGETARCH
COPY $BIN_DIR/$TARGETARCH/mirakurun_exporter /usr/bin/mirakurun_exporter

EXPOSE 9110
ENTRYPOINT ["/usr/bin/mirakurun_exporter"]
