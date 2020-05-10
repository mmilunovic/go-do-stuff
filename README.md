# Go Do Stuff

Go Do Stuff is a simple TODO App written in GoLang.

It contains two services:
* Worker - RESTful API that does basic CRUD operations
* Load Balancer - Simple load balancer that uses Round Robin algorithm to distribute requests to workers

SQLite is used for storage.
## Installation

```bash
git clone https://github.com/mmilunovic/go-do-stuff.git
```
#### Dependencies

```bash
go get github.com/jinzhu/gorm
go get github.com/gorilla/mux
go get github.com/jinzhu/gorm/dialects/sqlite
go get github.com/dgrijalva/jwt-go
```

## Usage

To start the Load Balancer do the following:
```bash
cd loadBalancer
go build
./loadBalancer
```

To start the worker:
``` bash
cd worker
go build
./worker
```