FROM golang:latest

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
                    git \
                    runit && \
    rm -rf /var/lib/apt/lists/*

ENV GOPATH /go
RUN mkdir -p $GOPATH
ENV PROJECT_NAME go.roman.zone/publisher
ENV PROJECT_PATH $GOPATH/$PROJECT_NAME

RUN go get -u github.com/peterbourgon/runsvinit
RUN cp $GOPATH/bin/runsvinit /usr/local/bin/
COPY ./publisher.service /etc/sv/publisher/run
RUN chmod 755 /etc/sv/publisher/run && \
    ln -sf /etc/sv/publisher /etc/service/

COPY . $PROJECT_PATH
WORKDIR $PROJECT_PATH

RUN go install $PROJECT_NAME
RUN cp $GOPATH/bin/publisher /usr/local/bin/

# Cleanup
RUN rm -rf $GOPATH

EXPOSE 80
ENTRYPOINT ["/usr/local/bin/runsvinit"]
