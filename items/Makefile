dev:
	docker-compose -f docker-compose.yml up

dev-db:
	docker run -d --name postgres_items_db -p 5432:5432 --restart always --env-file ./.env.dev -v pg_data_items:/var/lib/postgresql/data postgres:15.0-alpine