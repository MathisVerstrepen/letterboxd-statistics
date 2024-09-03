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
FROM alpine:3.14

# Set the working directory
WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main .
COPY --from=builder /app/assets ./assets

# Command to run the executable
CMD ["./main"]
