FROM golang:1.16.8-alpine3.14

WORKDIR /app

RUN apk update && apk add --no-cache gcc musl-dev git

COPY . .

RUN go build -ldflags '-w -s' -a -o ./bin ./cmd

CMD ["/app/bin"]
EXPOSE 9000
EXPOSE 7000