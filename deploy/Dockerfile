FROM golang:1.22.3-alpine3.19 as bg-builder
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
WORKDIR /
COPY --from=bg-builder /app/app-binary /app-binary
COPY --from=bg-builder /app/front /front
ENTRYPOINT [ "/app-binary", "-env", "production" ]
