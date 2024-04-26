FROM golang:1.22.0-alpine3.19 AS builder
LABEL authors="mirza.hilmi@gmail.com"
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o app ./cmd/api

FROM alpine:3.19.1
WORKDIR /app

ARG USER=runner
ARG GROUP=$USER

RUN addgroup -g 1000 runner && \
adduser -DH -g '' -G runner -u 1000 runner && \
apk add dumb-init

COPY --from=builder --chown=$USER:$GROUP --chmod=500 \
/src/app \
/src/.env \
/src/*.html \
./

USER $USER:$GROUP

ENTRYPOINT [ "/usr/bin/dumb-init", "--" ]
CMD [ "./app" ]
