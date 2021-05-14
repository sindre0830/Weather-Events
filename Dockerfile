FROM golang:1.16

LABEL maintainer "rickarl@stud.ntnu.no"

RUN mkdir /app
ADD . /app
WORKDIR /app

ADD ./go.mod /
ADD ./main.go /

RUN CGO_ENABLED=0 GOOS=linux go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o main

RUN chmod +x ./main

CMD ["./main"]
