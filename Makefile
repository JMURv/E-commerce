build_reviews:
	docker build . -t reviews_svc -f ./reviews/Dockerfile  --no-cache

consul:
	docker run --net=host -d -p 8500:8500 -p 8600:8600/udp --name=consul hashicorp/consul:latest agent -dev -ui -client="0.0.0.0"

kafka:
	docker run --net=host -d -p 9092:9092 --name=kafka apache/kafka:3.7.0

