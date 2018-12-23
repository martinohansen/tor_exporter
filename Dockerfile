FROM golang

ADD . /go/src/github.com/martinohansen/tor_exporter
RUN go install github.com/martinohansen/tor_exporter

ENTRYPOINT ["/go/bin/tor_exporter"]
CMD ["--tor.control-socket", "/var/run/tor/control"]

EXPOSE 9105
