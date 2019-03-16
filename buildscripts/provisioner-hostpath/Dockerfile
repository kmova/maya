FROM alpine:3.6

RUN apk add --no-cache \
    iproute2 \
    bash \
    curl \
    net-tools \
    mii-tool \
    procps \
    libc6-compat \
    ca-certificates

RUN mkdir -p /etc/apiserver/orchprovider

COPY provisioner-hostpath /

ARG BUILD_DATE

LABEL org.label-schema.name="provisioner-hostpath"
LABEL org.label-schema.description="API server for OpenEBS"
LABEL org.label-schema.url="http://www.openebs.io/"
LABEL org.label-schema.vcs-url="https://github.com/openebs/maya"
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.build-date=$BUILD_DATE

CMD ["/provisioner-hostpath"]
