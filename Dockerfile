##Latest golang base image
#FROM golang:1.15.8-alpine3.13
#
##Maintainer Info
#LABEL maintainer="Dmitry Fedorov <dimdomfedorov@gmail.com>"
#
##Working directory
#WORKDIR /cmd
#
##Copy the source from the current directory to the container workdir
#COPY . .
#
##Download all dependencies
##RUN go mod download
#
##Build the Go app
#RUN go build -o main /cmd/cmd/server/main.go
#
##Comand to run the executable
#CMD ["./main"]

FROM golang:1.15.8-alpine3.13 AS go_image
WORKDIR /go_cmd
COPY . .
RUN go build -o main /go_cmd/cmd/server/main.go

FROM alpine:latest
WORKDIR /alpine_cmd
COPY --from=go_image /go_cmd/main /alpine_cmd
CMD ["./main"]
