up:
    docker-compose up --build

run-app:
    ./run-app.sh

test: 
    migrate-up
    go test -v ./tests -cover

migrate-up:
    migrate -path datastore/postgres/migrations -database "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" -verbose up

migrate-down:
    migrate -path datastore/postgres/migrations -database "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" -verbose down --all