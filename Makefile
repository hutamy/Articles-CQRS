dep:
	go mod tidy && go mod vendor

up:
	docker compose up -d --build

down:
	docker compose down