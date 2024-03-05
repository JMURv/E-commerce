go:
	go run main.go

dev:
	docker-compose up

down:
	docker-compose down

build_reviews:
	docker build . -t reviews_srv:latest -f ./reviews/Dockerfile  --no-cache

run_reviews:
	docker run --net=host --rm --name=reviews reviews_srv:latest

run_discovery:
	docker run --net=host -d -p 8500:8500 -p 8600:8600/udp --name=consul hashicorp/consul:latest agent -dev -ui -client="0.0.0.0"
