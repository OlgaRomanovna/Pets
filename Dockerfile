FROM golang:1.24

WORKDIR /app

COPY . .

RUN go build -mod=vendor -o petfeed ./cmd/petfeed

EXPOSE 8080

CMD ["./petfeed"]
