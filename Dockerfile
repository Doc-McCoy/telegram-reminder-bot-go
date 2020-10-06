FROM golang

ADD ./src /go/src/bot

RUN go get -u gorm.io/gorm
RUN go get -u gorm.io/driver/postgres
RUN go get -u github.com/go-telegram-bot-api/telegram-bot-api

RUN go install bot

ENTRYPOINT /go/bin/bot
