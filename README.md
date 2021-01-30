# Atwell

Atwell is a Twitter for one person.

## Features
- Clean architecture
- OpenAPI(Swagger)
- CI
- Testing
- Docker

## How to use
### Set up
#### Database
Create a mysql database and schema by the following migration command.
```shell
cd ./migration
# migration
goose mysql "user:pass@(host:port)/db_name?parseTime=true" up

# confirm to apply sql
goose mysql "user:pass@(host:port)/db_name?parseTime=true" status
```

### Configuration
Please set following environment variables.

- ATWELL_DB_HOST
- ATWELL_DB_PORT
- ATWELL_DB_USER
- ATWELL_DB_PASSWORD
- ATWELL_DB_DBNAME

If you run test at local, set test db info.

- ATWELL_TEST_DB_HOST
- ATWELL_TEST_DB_PORT
- ATWELL_TEST_DB_USER
- ATWELL_TEST_DB_PASSWORD
- ATWELL_TEST_DB_DBNAME

### Start server
```shell
go run main.go
```

Or you can use docker integration. Please see [docker-compose file](./docker/docker-compose.yml).

## API Documentation
https://riku-1.github.io/atwell/

The documentation is created by [swag](https://github.com/swaggo/swag).  
Swag is a tool converting Go annotations to Swagger Documentation.

```shell
swag init --parseDependency --parseInternal
```

### Testing

## Using Stacks/Libraries

