# Start from a Debian image with the latest version of Go installed

FROM golang

RUN go get gopkg.in/DuoSoftware/DVP-DashBoard.v2/DashBoard

RUN go install gopkg.in/DuoSoftware/DVP-DashBoard.v2/DashBoard

ENTRYPOINT /go/bin/DashBoard

EXPOSE 8841
