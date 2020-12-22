FROM golang:1.15.6-buster

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go install

ENTRYPOINT [ "/go/bin/andaman" ]