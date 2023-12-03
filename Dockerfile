FROM golang:latest

WORKDIR /app

COPY . .

RUN go get -u github.com/gorilla/mux
RUN go get -u github.com/lib/pq
RUN go get -u github.com/joho/godotenv
RUN go get -u gorm.io/driver/postgres
RUN go get -u gorm.io/gorm