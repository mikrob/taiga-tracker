FROM golang:1.6

COPY id_rsa_jenkins /tmp/
RUN mkdir ~/.ssh/
RUN ssh-keyscan -H gitlab.botsunit.com >> ~/.ssh/known_hosts
RUN eval $(ssh-agent) && ssh-add /tmp/id_rsa_jenkins && git clone git@gitlab.botsunit.com:infra/taiga-gitlab.git /go/src/gitlab.botsunit.com/infra/taiga-gitlab
RUN rm /tmp/id_rsa_jenkins

COPY . /go/src/taiga-tracker
WORKDIR /go/src/taiga-tracker
EXPOSE 8282
RUN go get -v && go build
CMD ["taiga-tracker"]
