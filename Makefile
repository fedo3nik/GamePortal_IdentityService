d_build:
	docker build -t fedo3nik/game_portal_identity_service:1.0.0 .

d_push:
	docker push fedo3nik/game_portal_identity_service:1.0.0

d_run:
	docker run -p 8080:8080 --env-file=config.env fedo3nik/game_portal_identity_service:1.0.0

go_lint:
	docker run --rm -v ${PWD}:/app -w /app/ golangci/golangci-lint:v1.35.2 golangci-lint run -v --timeout=5m

mongo_run:
	docker run -it -v mongodata:/data/db -p 27017:27017 --name mongodb -d mongo

cover_profile:
	go test -coverprofile coverage.out ./...

cover_html:
	go tool cover -html=coverage.out