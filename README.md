
# Project-Ecommerce

This project is written purely in Go and used PostgreSQL for its database needs
.Entirely Dockerized this project using dockerfile and docker-compose.yml

### Framworks used
Gin-Gonic:This whole project is fully completed using gin-gonic which is a popular go framework for rapid web development
```
go get -u github.com/gin-gonic/gin
```

### Database used:
PostgreSQL:This project mainly used PostgreSQL as Database with the help of ORM tool named GORM.It provides better and simplified forms of queries for better understanding

```
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
```
### UP & RUN Project
Run this project in your local repository using Docker(if Docker installed in your system) or Go run command

commands to run using docker:
```
docker compose up -d
docker compose down
```
commands to run using go run:
```
go run main.go
```

### Use API Platform
API platforms such as Postman can be used to run all the API's Provided by this project

### API Documentation
```
https://documenter.getpostman.com/view/23476254/2s8YzL46Xq
```
