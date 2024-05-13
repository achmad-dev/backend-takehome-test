FROM golang:1.18-buster  AS builder

ENV GO111MODULE=on 

WORKDIR /opt/backend-recruit
RUN apt-get -y update
RUN apt-get install -y gcc musl-dev git ca-certificates
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY . .

ENV TZ Asia/Jakarta

RUN --mount=type=cache,target=/root/.cache/go-build go build ./cmd/app/main.go


FROM bitnami/minideb:latest
WORKDIR /opt/backend-recruit
RUN mkdir config
COPY --from=builder /opt/backend-recruit/main /opt/backend-recruit/main
ADD ./config/app.env.docker ./config/app.env
RUN chmod +x /opt/backend-recruit/main