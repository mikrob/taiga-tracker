FROM golang:1.6
COPY . /go/src/taiga_tracker
WORKDIR /go/src/taiga_tracker
EXPOSE 8080
RUN go get -v && go build
CMD ["taiga-tracker"]
