FROM golang:1.18.3-alpine as builder
RUN apk add ca-certificates git
ARG current_time
ARG git_description
COPY ./ /app
WORKDIR /app
RUN CGO_ENABLED=0 go build -ldflags \
    "-s -X main.buildTime=${current_time} -X main.version=${git_description}" \
    -o ./app-binary ./cmd && \
    mv ./app-binary /app/ && \
    chmod +x /app/app-binary

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder /app/app-binary /app-binary
ENTRYPOINT [ "/app-binary", "-env", "production" ]
