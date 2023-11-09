FROM golang:1.19-bullseye AS builder
ARG GITHUB_TOKEN
ARG GONOPROXY=github.com/farawaygg/*
ARG GOPRIVATE=github.com/farawaygg/*
ARG GONOSUMDB=github.com/farawaygg
ARG GOGC=off
WORKDIR /app
COPY . ./
RUN git config --global url."https://${GITHUB_TOKEN}:@github.com/".insteadOf https://github.com/
RUN go build -o /binary './cmd/auth'

FROM gcr.io/distroless/base-debian11:nonroot
USER nonroot:nonroot
WORKDIR /
COPY --from=builder /binary /binary
CMD ["/binary", "-config", "/etc/app/config.yaml"]
