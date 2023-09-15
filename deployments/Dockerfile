FROM golang:1.17-alpine AS build

ARG CGO_ENABLED=0
ARG GO111MODULE=on
ARG GOPRIVATE=git.cyradar.com
ARG GIT_USERNAME
ARG GIT_ACCESS_TOKEN

WORKDIR /go/src/app

RUN apk add git

RUN echo https://${GIT_USERNAME}:${GIT_ACCESS_TOKEN}@git.cyradar.com > ~/.git-credentials
RUN echo -e "[credential]\n\thelper = store" > ~/.gitconfig

COPY go.mod .
COPY go.sum .

RUN go mod download -x

COPY . .

RUN go install -v ./...

FROM alpine

WORKDIR /app
COPY --from=build /go/bin/caddy .
RUN mkdir -p configs && touch configs/.env && apk --no-cache add tzdata

ENTRYPOINT ["/app/caddy"]
# CMD ["--addr", "0.0.0.0:80"]

# ENTRYPOINT ["/app/endpoint-security-server", "--env", "/app/configs/.env", "serve"]
# CMD ["--addr", "0.0.0.0:80"]

# docker build --build-arg GIT_USERNAME --build-arg GIT_ACCESS_TOKEN -t cyradardev/ess -f deployments/Dockerfile .
