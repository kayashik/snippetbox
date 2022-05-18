# syntax=docker/dockerfile:1

FROM golang:1.17-alpine
WORKDIR /app
COPY . ./
RUN go mod download
RUN go build -o /snippetbox  ./cmd/web
EXPOSE 4000
CMD [ "/snippetbox" ]
