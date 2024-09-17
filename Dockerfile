FROM golang:latest

WORKDIR /go/src

COPY ./src /go/src

RUN go get -u gorm.io/gorm
RUN go get -u gorm.io/driver/mysql
RUN go get -u github.com/labstack/echo/v4
RUN go mod tidy

EXPOSE 6050