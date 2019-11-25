FROM sc4h/alpine-s6overlay:3.10

ENV \
  CONFIG_DIR=/config/trackarr

# Binary
COPY ["dist/trackarr_linux_amd64/trackarr", "/app/trackarr/trackarr"]

# Add root files
COPY ["build/root/", "/"]

VOLUME ["/config"]
