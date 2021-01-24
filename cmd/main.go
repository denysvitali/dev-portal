package main

import (
	"github.com/denysvitali/dev-portal/pkg/server"
	"github.com/denysvitali/dev-portal/pkg/server/app"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
)

func main(){
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	
	file, err := os.Open("./config.yml")
	if err != nil {
		logger.Fatal(err)
	}
	var config app.Config
	
	d := yaml.NewDecoder(file)
	err = d.Decode(&config)
	
	if err != nil {
		logger.Fatal(err)
	}

	s := server.New(config.ListenAddr, logger, &config)
	s.Start()
}