#Latest golang base image
FROM golang:1.15

#Maintainer Info
LABEL maintainer="Dmitry Fedorov <dimdomfedorov@gmail.com>"

#Working directory
WORKDIR /cmd

#Copy the source from the current directory to the container workdir
COPY . .

#Download all dependencies
RUN go mod download

#Build the Go app
RUN go build -o main /cmd/cmd/main.go

#Expose port 8080 to outside world
EXPOSE 8080

#Comand to run the executable
CMD ["./main"]
