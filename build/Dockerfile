FROM sc4h/alpine-s6overlay:3.11

ENV \
  CONFIG_DIR=/config

# Binary
COPY ["dist/trackarr_linux_amd64/trackarr", "/app/trackarr/trackarr"]

# Add root files
COPY ["build/root/", "/"]

EXPOSE 7337

VOLUME ["/config"]
