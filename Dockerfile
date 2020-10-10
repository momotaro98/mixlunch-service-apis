FROM golang:1.13.1-stretch as builder

WORKDIR /github.com/momotaro98/mixlunch-service-api

ENV GO111MODULE="on"
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o mixlunch-service-api

FROM alpine:3.8

WORKDIR /root/

# Need CA certificates to use HTTPS in a container
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

# Multi stage build function of Docker
COPY --from=builder /github.com/momotaro98/mixlunch-service-api/mixlunch-service-api .
COPY --from=builder /github.com/momotaro98/mixlunch-service-api/serviceAccount/serviceAccountKey.json ./serviceAccount/

CMD ["./mixlunch-service-api"]
