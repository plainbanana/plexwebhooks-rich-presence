# lapis builder image
FROM golang:latest as builder
LABEL maintainer "plainbanana <kazukidegozaimasuruzo@gmail.com>"
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /
COPY . .
RUN go build -o app

# lapis image
# docker run ${containername} --env-file .env
FROM alpine
LABEL maintainer "plainbanana <kazukidegozaimasuruzo@gmail.com>"
ENV DOTENV=false
RUN apk add --no-cache ca-certificates
COPY --from=builder /app /
CMD ["/app"]