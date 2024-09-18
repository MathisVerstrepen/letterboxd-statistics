# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make npm

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

# Set the working directory
WORKDIR /app

# Install npm dependencies
COPY package.json ./
RUN npm install

# Copy go mod and sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy && go mod download

# Copy the source code and generate the templ files
COPY . .
RUN templ generate

# Build the CSS
RUN npx tailwindcss -i ./assets/css/input.css -o ./assets/css/output.css --minify

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main main.go

# Download the assets
RUN wget https://raw.githubusercontent.com/MathisVerstrepen/templ_ui_library/main/assets/css/output.css -O ./assets/css/templ_ui_library.css

# Final stage
FROM debian:buster-slim

# Install phantomjs
# from https://github.com/wernight/docker-phantomjs/blob/master/latest/Dockerfile
ENV PHANTOM_JS_VERSION ${PHANTOM_JS_VERSION:-2.1.1-linux-x86_64}

# Install runtime dependencies
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
            ca-certificates \
            curl \
            bzip2 \
            libfontconfig \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

ENV DUMB_INIT_VERSION ${DUMB_INIT_VERSION:-1.2.2}

RUN set -x  \
    # Install official PhantomJS release
    && mkdir /tmp/phantomjs \
    && curl -Ls https://bitbucket.org/ariya/phantomjs/downloads/phantomjs-${PHANTOM_JS_VERSION}.tar.bz2 \
            | tar -xj --strip-components=1 -C /tmp/phantomjs \
    && mv /tmp/phantomjs/bin/phantomjs /usr/local/bin

ENV OPENSSL_CONF=/opt/openssl.cnf

RUN set -x \
    # Runs as non-root user.
    && useradd --system --uid 52379 -m --shell /usr/sbin/nologin phantomjs \
    && touch /opt/openssl.cnf

# Set the working directory
WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main .
COPY --from=builder /app/assets ./assets

# Command to run the executable
CMD ["./main"]
