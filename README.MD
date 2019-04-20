## download && install dependencies
```cmd
go get -u github.com/jinzhu/gorm
go get -u github.com/gin-gonic/gin
go get -u github.com/go-sql-driver/mysql
go get -u github.com/gomodule/redigo/redis
```

## build
```cmd
# set GOARCH and GOOS

# build at root folder of project
go build -i -o sugar_server main.go 
```
## parameters for executing
```
  -gen
        generate a configuration example file.
  -help
        cat help information
  -init
        init the database
  -port uint
        set port of server (default 8080)
  -release
        run server on release
  -server
        run server on debug

```
## sequence for executing
`./sugar_server -help`  
 
`./sugar_server -gen` generate a configuration example file.

`./sugar_server -init` init data of database

`./sugar_server -server -release -port 8000` run production server on port 8000 