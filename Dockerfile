FROM docker.io/golang:bookworm 

WORKDIR /maobot

COPY . .

RUN env GOOS=linux GOARCH=arm64 go build -o maobot

CMD [ "./maobot" ]

