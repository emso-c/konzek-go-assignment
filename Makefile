.PHONY: build-and-run

build-and-run: docker-compose up -d --build

run:
	docker-compose up -d

build:
	docker-compose build

follow-logs:
	docker-compose logs -f

stop:
	docker-compose down

check-db:
	docker exec go_assignment-database-1 pg_isready

run-tests:
	cd app && docker build -f Dockerfile.test -t go_assignment_test . && docker run go_assignment_test