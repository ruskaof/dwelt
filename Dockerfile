FROM golang:1.22

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o main

EXPOSE 8080

ARG DWELT_JWT_KEY
ENV DWELT_JWT_KEY=$DWELT_JWT_KEY
ARG DWELT_WORKFLOW_RUN_NUMBER
ENV DWELT_WORKFLOW_RUN_NUMBER=$DWELT_WORKFLOW_RUN_NUMBER

ARG DB_HOST
ENV DB_HOST=$DB_HOST
ARG DB_PORT
ENV DB_PORT=$DB_PORT
ARG DB_USER
ENV DB_USER=$DB_USER
ARG DB_PASSWORD
ENV DB_PASSWORD=$DB_PASSWORD
ARG DB_NAME
ENV DB_NAME=$DB_NAME


CMD ./main
