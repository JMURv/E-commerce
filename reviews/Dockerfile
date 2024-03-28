FROM golang:1.21.5-alpine3.19

WORKDIR /app

COPY ../go.mod ../go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/reviews/cmd
CMD ["go", "run", "."]