run-local:
	nodemon -e go -x 'go run . || exit 1' --signal SIGTERM

build:
	docker build -t todo .

run-docker:
	docker run -d -p 8080:8080 --name todo-app todo
