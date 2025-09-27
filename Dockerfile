FROM golang:1.25.1-alpine

WORKDIR /app

RUN apk add curl jq

RUN curl -s https://api.github.com/repos/cmiguelrb/benthos/releases/latest \
  | jq -r '.assets[] | select(.name | test("benthos-linux-arm64")) | .browser_download_url' \
  | xargs curl -L -o benthos && chmod +x benthos

##Env variables in the docker_compose file

EXPOSE 3800

CMD ["./benthos"]