FROM golang

ADD . /go/src/reminder-bot

RUN go get gorm.io/gorm
RUN go get gorm.io/driver/postgres
RUN go get github.com/go-telegram-bot-api/telegram-bot-api@develop
RUN go install reminder-bot

ENTRYPOINT /go/bin/reminder-bot

EXPOSE 8080
