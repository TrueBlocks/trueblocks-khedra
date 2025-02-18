# syntax=docker/dockerfile:1

# =========== 1) BUILD STAGE =============
FROM ubuntu:24.04 AS builder

# Install dependencies needed to build trueblocks-khedra
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
      build-essential \
      cmake \
      git

# Create and change to the source directory
WORKDIR /app

# Copy source into the container
COPY . /app

# Build the project
RUN mkdir build && cd build && cmake .. && make

# =========== 2) FINAL STAGE =============
FROM ubuntu:22.04

# Copy compiled binary from the builder stage
COPY --from=builder /app/build/khedra /usr/local/bin/khedra

# Copy example config (you can move or rename as preferred)
COPY --from=builder /app/config.example /root/.trueblocks/trueblocks-khedra.conf

# Set environment variables or defaults for Khedra
ENV KHEDRA_CONFIG=/root/.trueblocks/trueblocks-khedra.conf \
    KHEDRA_DATA_DIR=/root/.trueblocks/data \
    KHEDRA_LOG_LEVEL=INFO

# Default entrypoint runs 'khedra' 
ENTRYPOINT ["khedra"]
# Default command shows help text
CMD ["--help"]
