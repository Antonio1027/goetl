SHELL := /bin/bash

run:
	docker-compose up --build


stop:
	docker-compose down


restart: stop run


tests:
	docker build -f Dockerfile.multistage -t goetl-test --progress plain --no-cache --target run-test-stage .