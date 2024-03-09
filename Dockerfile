FROM golang:1.22

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o main

EXPOSE 8080
ARG JWT_KEY

CMD ./main -jwtkey "$JWT_KEY"
