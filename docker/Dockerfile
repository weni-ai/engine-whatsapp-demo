FROM golang:1.16.8-bullseye AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download -x

COPY . ./

RUN GOOS=linux GOARCH=amd64 go build -o /bin/engine-whatsapp-demo ./cmd/main.go

FROM gcr.io/distroless/base-debian11

WORKDIR /app

COPY --from=builder bin/engine-whatsapp-demo ./
EXPOSE 9000
EXPOSE 7000
CMD ["./engine-whatsapp-demo"]