FROM golang:1.25
WORKDIR /server
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o s ./cmd/server.go
CMD [ "./s" ]
