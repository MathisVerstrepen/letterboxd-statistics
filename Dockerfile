FROM golang:1.23-bookworm

WORKDIR /app

# Install cron
RUN apt-get update && apt-get -y install cron
COPY crontab /etc/cron.d/crontab
RUN chmod 0644 /etc/cron.d/crontab && /usr/bin/crontab /etc/cron.d/crontab

COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN go build -o /app/main .

RUN chmod +x /app/entrypoint.sh