FROM golang:1.23 AS builder

WORKDIR /app

COPY . ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /cyoa

FROM alpine

COPY --from=builder /cyoa /cyoa
COPY ./templates /templates/
COPY ./stories /stories/

EXPOSE 8080

CMD ["/cyoa", "-t", "http", "-f", "./stories/story.json"]