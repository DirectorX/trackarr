FROM sc4h/alpine-s6overlay:3.12

ARG ARCH

ENV \
  CONFIG_DIR=/config

# Binary
COPY ["dist/trackarr_linux_${ARCH}/trackarr", "/app/trackarr/trackarr"]

# Add root files
COPY ["build/root/", "/"]

EXPOSE 7337

VOLUME ["/config"]
