package main

import (
	"fmt"
	"log"
	"net/http"

	"screen-shot-api/api"
	"screen-shot-api/config"
	"screen-shot-api/db"
	"screen-shot-api/logger"

	_ "github.com/lib/pq"
)

func main() {
	/* configuration initialize start */
	c := config.Configuration()
	/* configuration initialize end */

	/* logger initialize start */
	mylogger := logger.NewZapLogger()
	logger.InitializeLogger(&mylogger)
	defer logger.Close()
	/* logger initialize end */

	/* initialize database start */
	if err := db.InitDB(c.DBEngine, c.DBConnectionString); err != nil {
		panic(err)
	}
	defer db.Close()
	/* initialize database end */

	/* initialize webserver start */
	r := api.NewRouter()
	addr := fmt.Sprintf("%v:%d", c.WebAddress, c.WebPort)
	log.Fatal(http.ListenAndServe(addr, r))
	/* initialize webserver end */
}
