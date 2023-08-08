FROM golang:alpine as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o yt-dlp-archive-viewer

FROM alpine:3.18.3

COPY --from=builder /app/yt-dlp-archive-viewer /usr/bin

EXPOSE 8000

CMD [ "/usr/bin/yt-dlp-archive-viewer" ]
