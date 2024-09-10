# Dockerfile
FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /go-todo-list

EXPOSE 8000

CMD [ "/go-todo-list" ]
