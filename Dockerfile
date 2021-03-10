FROM golang:1.15.8-alpine3.13 AS go_image
WORKDIR /go_cmd
COPY . .
RUN go build -o main /go_cmd/cmd/server/main.go

FROM alpine:latest
WORKDIR /alpine_cmd
COPY --from=go_image /go_cmd/main /alpine_cmd
EXPOSE 8080
CMD ["./main"]
