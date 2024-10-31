FROM golang:1.22-alpine as builder

WORKDIR /source
COPY . /source

RUN go mod download
RUN go build -v -o prometheus-dummy-exporter .

FROM        alpine:latest
MAINTAINER  Dongho Jung <dongho971220@gmail.com>

COPY --from=builder /source/prometheus-dummy-exporter  /bin/prometheus-dummy-exporter

EXPOSE     2112
ENTRYPOINT [ "/bin/prometheus-dummy-exporter" ]
