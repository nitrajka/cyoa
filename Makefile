build:
	docker build --tag cyoa:1.0 .
run:
	make build && docker run -p 8080:8080 cyoa:1.0