# Start from a Debian image with the latest version of Go installed

FROM golang
ARG MAJOR_VER

RUN go get gopkg.in/DuoSoftware/DVP-DashBoard.$MAJOR_VER/DashBoard

RUN go install gopkg.in/DuoSoftware/DVP-DashBoard.$MAJOR_VER/DashBoard

ENTRYPOINT /go/bin/DashBoard

EXPOSE 8841
