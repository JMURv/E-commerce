build_reviews:
	docker build . -t reviews_svc -f ./reviews/Dockerfile  --no-cache

consul:
	docker run d -p 8500:8500 -p 8600:8600/udp --name=consul hashicorp/consul:latest agent -dev -ui -client="0.0.0.0"

kafka:
	docker run -d -p 9092:9092 --name=kafka apache/kafka:3.7.0

redis:
	docker run -d -p 6379:6379 --name=redis redis:alpine