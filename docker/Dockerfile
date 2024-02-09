FROM golang

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o bin ./cmd

EXPOSE 8080

ENTRYPOINT ["./bin"]
