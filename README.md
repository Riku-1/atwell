# Atwell

Atwell is a Twitter for one person.

## Features
- Clean architecture
- OpenAPI(Swagger)
- CI
- Testing
- Docker

## API Documentation
https://riku-1.github.io/atwell/

The documentation is created by [swag](https://github.com/swaggo/swag).  
Swag is a tool converting Go annotations to Swagger Documentation.

```shell
swag init --parseDependency --parseInternal
```

## How to use
### Set up
#### Database
TODO: Create migration description

Set up docker-compose environment

### Start server
```shell
docker-compose up -d
```

## Using Stacks/Libraries
|Stack|Description|
|---|---|
|[Gorm](https://gorm.io/)|Golang ORM|
|gin||
|goose|Golang ORM|
