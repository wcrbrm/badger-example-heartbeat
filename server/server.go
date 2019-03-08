package main

import (
	_ "github.com/dgraph-io/badger"
	_ "github.com/gin-gonic/gin"
	_ "github.com/Depado/ginprom"
	_ "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"fmt"
)

func main() {
	fmt.Println("Started")
}
