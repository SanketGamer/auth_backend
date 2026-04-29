## Manual installation
 - Install Golang locally
 - Clone the repo
 - Install dependencies (`go mod download`)
 - Start the DB locally
    - `docker run -e POSTGRES_PASSWORD=mysecretpassword -e POSTGRES_USER=postgres -e POSTGRES_DB=authdb -d -p 5432:5432 postgres`
    - Go to neon.tech and get yourself a new DB
 - Change the `.env` file and update your DB credentials
 - Build the project (`go build -o main ./cmd`)
 - Start the server (`./main`)

## Docker installation
 - Install docker
 - Create a network - `docker network create auth_project`
 - Start postgres
    - `docker run --network auth_project --name postgres_db -e POSTGRES_PASSWORD=mysecretpassword -e POSTGRES_USER=postgres -e POSTGRES_DB=authdb -d -p 5432:5432 postgres`
 - Build the image - `docker build -t auth-jwt-service .`
 - Start the image - `docker run -e DB_HOST=postgres_db -e DB_PORT=5432 -e DB_USER=postgres -e DB_PASSWORD=mysecretpassword -e DB_NAME=authdb -e DB_SSLMODE=disable --network auth_project -p 8080:8080 auth-jwt-service`

### to create 2 container one go-app and postgres container and use volume and connect via network
- postgres - `docker run --name postgres_db --network auth_project -v postgres_data:/var/lib/postgresql/data -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=mysecretpassword -e POSTGRES_DB=authdb -d postgres:14`
- Go-App - `docker run --name go_api --network auth_project -p 8080:8080 -e DB_HOST=postgres_db -e DB_PORT=5432 -e DB_USER=postgres -e DB_PASS=mysecretpassword -e DB_NAME=authdb -e DB_SSL=disable auth-jwt-service`


## Docker Compose installation steps
 - Install docker, docker-compose
 - Run `docker-compose up`  new one is `docker compose up --build`