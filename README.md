# Overview
It's an API Fiber Starter.
- Minimum Go Requirement 1.15
- Fiber framework [link](https://gofiber.io)

## What's inside:
- Environment configuration
- Migrations
- Authentication with JWT

## Usage
1. COPY .env.example TO .env
    ``` ~ cp -r .env.example .env ```
2. Generate APP Key
    Please visit http://www.sha1-online.com to generate a new key and add the generated key to .env
    ``` APP_KEY=COPY HERE ```
3. Install all dependencies
    ```~ go mod download```
4. Migrations
    Install the migration tool and follow this [link](https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md) for the instructions.
    **Create Migration**
    please follow this [link](https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md)
    example: 
    ``` migrate create -ext sql -dir db/migration -seq create_[table name]_table ```
    **Execute Migration**
    add an env variable for psql connection to your system:
    ```export POSTGRESQL_URL='postgres://user:pass@localhost:5432/dbname?sslmode=disable&search_path=public'```
    and run the command to execute migrations:
    ``` migrate -database ${POSTGRESQL_URL} -path db/migration up ```
5. Run application using the command in the terminal:
    `go run main.go api` OR `go run main.go cmd [flag]` 
6. Run cli application for queue using the command in the terminal:
    `go run main.go cmd -queue=queue_name` 

## License
The project is developed by [Devrian]
