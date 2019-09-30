FROM golang:1.12-alpine as build

RUN apk add --no-cache git

COPY go.mod go.sum /app/
RUN cd /app && go mod download

COPY . /app

RUN cd /app && go build -o service-manager cmd/servicemanager/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=build /app/service-manager /home/service-manager

RUN chmod +x /home/service-manager

EXPOSE 8080

CMD ["/home/service-manager"]