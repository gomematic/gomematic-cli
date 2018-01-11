FROM webhippie/alpine:latest

LABEL maintainer="Thomas Boerger <thomas@webhippie.de>" \
  org.label-schema.name="Gomematic CLI" \
  org.label-schema.vendor="Thomas Boerger" \
  org.label-schema.schema-version="1.0"

ENTRYPOINT ["/usr/bin/gomematic-cli"]
CMD ["help"]

RUN apk add --no-cache ca-certificates mailcap bash

COPY dist/binaries/gomematic-cli-*-linux-amd64 /usr/bin/
