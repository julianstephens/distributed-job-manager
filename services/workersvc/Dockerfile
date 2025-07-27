FROM golang:1.24-bookworm as base

# Load build arg
ARG APP_HOME

# Create development stage and set app location in container to /app
FROM base as development
WORKDIR $APP_HOME

# Add bash shell and gcc toolkit
# RUN apt update && apt install build-base

# Install isolate
RUN echo 'deb [arch=amd64 signed-by=/etc/apt/keyrings/isolate.asc] http://www.ucw.cz/isolate/debian/ bookworm-isolate main' > /etc/apt/sources.list.d/non-free.list
RUN curl https://www.ucw.cz/isolate/debian/signing-key.asc >/etc/apt/keyrings/isolate.asc
RUN apt update && apt install isolate

# Install air
RUN go install github.com/air-verse/air@latest

# Copy go.mod and install go dependencies
COPY go.mod ./
RUN go mod download

# Copy source files
COPY . .

# Build and run binary with live reloading
CMD ["air", "-c", "services/workersvc/.air.toml"]
