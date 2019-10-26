FROM golang as builder
COPY main.go /go/src/github.com/shreddedbacon/ipcam-alarmserver/
COPY go.mod /go/src/github.com/shreddedbacon/ipcam-alarmserver/
COPY go.sum /go/src/github.com/shreddedbacon/ipcam-alarmserver/

WORKDIR /go/src/github.com/shreddedbacon/ipcam-alarmserver/
RUN CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o alarmserver .

FROM alpine
EXPOSE 15002
WORKDIR /app
COPY --from=builder /go/src/github.com/shreddedbacon/ipcam-alarmserver/ .
ENTRYPOINT [ "/app/alarmserver" ]