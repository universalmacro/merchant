FROM golang:1.21 as build-stage
WORKDIR /app
ARG GIT_TOKEN
ARG GIT_NAME
COPY go.mod go.sum ./
RUN go env -w GOPRIVATE=github.com/universalmacro/*
RUN git config --global url."https://${GIT_NAME}:${GIT_TOKEN}@github.com".insteadOf "https://github.com"
RUN go mod download
COPY . .
RUN go build -o /main
EXPOSE 8080
CMD [ "/main" ]