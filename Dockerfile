# Start from a Debian image with the latest version of Go installed

FROM golang

RUN go get gopkg.in/DuoSoftware/DVP-DashBoard.v2

RUN go install gopkg.in/DuoSoftware/DVP-DashBoard.v2

ENTRYPOINT /go/bin/DVP-DashBoard

EXPOSE 8841
