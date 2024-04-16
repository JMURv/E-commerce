build_reviews:
	docker build . -t reviews_svc -f ./reviews/Dockerfile  --no-cache

consul:
	docker run --rm --name=consul -p 8500:8500 -p 8600:8600/udp hashicorp/consul:latest agent -dev -ui -client="0.0.0.0"

kafka:
	docker run --rm --name=kafka -p 9092:9092 apache/kafka:3.7.0

redis:
	docker run --rm --name=redis -p 6379:6379 redis:alpine

jaeger:
	docker run --rm --name jaeger \
      -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
      -p 6831:6831/udp \
      -p 6832:6832/udp \
      -p 5778:5778 \
      -p 16686:16686 \
      -p 4317:4317 \
      -p 4318:4318 \
      -p 14250:14250 \
      -p 14268:14268 \
      -p 14269:14269 \
      -p 9411:9411 \
      jaegertracing/all-in-one:latest

prometheus:
	docker run --rm --name prometheus \
		-p 9090:9090 \
		-v $PWD/prometheus.yml:/etc/prometheus/prometheus.yml \
		-v prometheus-data:/prometheus \
 		prom/prometheus:latest
