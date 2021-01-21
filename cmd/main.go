package main

import (
	"github.com/denysvitali/dev-portal/pkg/server"
	"github.com/sirupsen/logrus"
)

func main(){
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	s := server.New(":8081", logger)
	s.Start()
}