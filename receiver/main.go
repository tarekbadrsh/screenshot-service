package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"receiver/config"
	"receiver/logger"
	"receiver/messaging"
	"receiver/service"
	"strings"

	"net/http"
)

func main() {
	/* configuration initialize start */
	c := config.Configuration()
	/* configuration initialize end */

	/* logger initialize start */

	// disable global log package Ref.https://stackoverflow.com/questions/10571182/how-to-disable-a-log-logger
	log.SetOutput(ioutil.Discard)

	mylogger := logger.NewZapLogger()
	logger.InitializeLogger(&mylogger)
	defer logger.Close()
	/* logger initialize end */

	/* web initialize start */
	if c.APIPort < 1 {
		logger.Fatalf("Invalid value. The value should be a positive integer (greater than zero), the value : %v", c.APIPort)
	}
	/* web initialize end */

	/* kafka initialize start */
	bs := strings.Split(c.KafkaBrokers, ",")
	kafkaClinet, err := messaging.NewSaramaSyncProducer(bs)
	if err != nil {
		logger.Fatal(err)
	}
	messaging.InitializeKafka(kafkaClinet)
	/* kafka initialize end */

	/* create a new *router instance */
	router := service.NewRouter()
	strport := fmt.Sprintf(":%v", c.APIPort)
	logger.Infof("Start Listen And Serve on %v", strport)

	if err := http.ListenAndServe(strport, router); err != nil {
		logger.Fatal(err)
	}
}
