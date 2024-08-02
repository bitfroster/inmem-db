# Inmem database
This is a simple implmentation of in memory database with support of transactions. It's implemented in Go.

## Requirements
- Go 1.18+

 or
- Docker

**Note:** As it's library, not a binary, I decided not to use minimal image technics like a scratch to reduce an image size.

## Usage
Code is implemented as a library, not as a standalone application. You can run the tests with one of the following commands:
```
go test -v
```
or
```
docker build -t inmemdb .
docker run -ti inmemdb
go test -v
```

Also you can add ```-cover``` flag to check the coverage.