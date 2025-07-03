package main

import (
	"context"
	"time"

	"github.com/VarthanV/golang-exercises-concurrency/fanout"
	"github.com/VarthanV/golang-exercises-concurrency/generators"
	"github.com/sirupsen/logrus"
)

func main() {
	generators.SimpleGeneratorDriver()
	generators.SimpleGeneratorFromFuncDriver()

	fanout.NaiveFanout()
	//fanout.FanOutWithSem()

	logrus.Info("Downloader driver")
	fanoutDownloader()
}

func fanoutDownloader() {
	d := fanout.NewDownloader(1)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()
	d.Driver(ctx)
}
