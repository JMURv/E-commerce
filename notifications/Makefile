dev:
	docker-compose -f docker-compose.yml up

dev-db:
	docker run -d --name postgres_notifications_db -p 5432:5432 --restart always --env-file ./.env.dev -v pg_data_notifications:/var/lib/postgresql/data postgres:15.0-alpine