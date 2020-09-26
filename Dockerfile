FROM golang:1.15-alpine

WORKDIR /usr/src/app
COPY . .

# Install sqlite3
RUN apk --update upgrade && \
    apk add sqlite g++ && \
    rm -rf /var/cache/apk/*

RUN go get .
RUN go build -o leader-board .
RUN go build -o seed db/seed/seed.go

RUN apk del g++

CMD ["/usr/src/app/leader-board"]

ENTRYPOINT ["./entrypoint.sh"]
