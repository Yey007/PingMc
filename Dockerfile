FROM golang:alpine

ENV CONTAINER=TRUE

WORKDIR /go/src/pingmc
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build

CMD ["./pingmc"]