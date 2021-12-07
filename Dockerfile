FROM golang as build

COPY . /project

WORKDIR /project

RUN go build -o /bin/highloadsrv -v ./cmd/

EXPOSE 8080

CMD highloadsrv

