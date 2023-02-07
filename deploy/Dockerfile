FROM golang:1.20.0-alpine3.17 as builder
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

FROM node:lts-alpine as fg-builder
WORKDIR /app
COPY ./front/package*.json /app/
RUN npm install 
COPY ./front/ /app
RUN npm run build

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder /app/app-binary /app-binary
ENTRYPOINT [ "/app-binary", "-env", "production" ]
