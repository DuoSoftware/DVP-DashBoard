# Start from a Debian image with the latest version of Go installed

FROM golang

RUN go get github.com/DuoSoftware/DVP-DashBoard

RUN go install github.com/DuoSoftware/DVP-DashBoard

ENTRYPOINT /go/bin/DVP-DashBoard

EXPOSE 8841