FROM golang:1.16

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./... && \
go install -v ./... && \
make build && \
cp application /go/bin

CMD ["application"]