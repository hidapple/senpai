FROM golang:1.8

WORKDIR /go/src/app
ENV SLACK_API_TOKEN "SLACK_API_TOKEN"
COPY *.go ./

RUN go-wrapper download
RUN go-wrapper install

CMD ["go-wrapper", "run"]
