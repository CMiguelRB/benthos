FROM golang:1.25.1-alpine

WORKDIR /app

RUN apk add curl jq

RUN curl -s https://api.github.com/repos/cmiguelrb/benthos/releases/latest \
  | jq -r '.assets[] | select(.name | test("benthos-linux-arm64")) | .browser_download_url' \
  | xargs curl -L -o benthos && chmod +x benthos

ENV ENV=xxx
## server
ENV PORT=xxx
## database
ENV DB_USER=xxx
ENV DB_PASSWORD=xxx
ENV DB_HOSTNAME=xxx
ENV DB_PORT=xxx
ENV DB_NAME=xxx
## security
ENV ENCRYPTION_KEY=xxx

EXPOSE 3800

CMD ["./benthos"]