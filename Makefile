d_build:
	docker build -t fedo3nik/game_portal_identity_service:1.0.0 .

d_push:
	docker push fedo3nik/game_portal_identity_service:1.0.0

d_run:
	docker run -p 8080:8080 --env-file=configDev.env fedo3nik/game_portal_identity_service:1.0.0

go_lint:
	docker run --rm -v ${PWD}:/app -w /app/ golangci/golangci-lint:v1.36-alpine golangci-lint run -v --timeout=5m

mongo_run:
	docker run -it -v mongodata:/data/db -p 27017:27017 --name mongodb -d mongo:4.4.4

cover_profile:
	go test -coverprofile coverage.out ./...

cover_html:
	go tool cover -html=coverage.out

ci_lint:
	circleci local execute --job lint

ci_test:
	circleci local execute --job test

ci_build:
	circleci local execute --job build_and_push_image -e DOCKER_REGISTRY_HOST=281520863489.dkr.ecr.eu-central-1.amazonaws.com -e STAGING_AWS_ACCESS_KEY_ID=AKIAI44U3T4U4TSNTP3Q -e STAGING_AWS_SECRET_ACCESS_KEY=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 -e STAGING_AWS_REGION=eu-central-1
