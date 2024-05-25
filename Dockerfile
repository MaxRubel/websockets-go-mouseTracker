FROM golang:1.22-alpine

WORKDIR /app
COPY WebsocketsGo .

EXPOSE 8069

CMD ["./WebsocketsGo"]
