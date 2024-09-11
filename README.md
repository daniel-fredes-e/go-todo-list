# go-todo-list

Task system for users with secure session

commands:

app:

docker-compose up --build

test:

docker build -t go-todo-list-tests -f Dockerfile.test .

sudo docker run --rm go-todo-list-tests

swagger:

http://localhost:4000/swagger/index.html#/
