package main

import (
	"os"

	"github.com/sorenhq/go-plugin-sdk/logtool"
)

const ServiceName = "soren-plugin-sdk-lab"

func init() {
	if os.Getenv("ENV") == "development" {
		logtool.Init(ServiceName, true)
	} else {
		logtool.Init(ServiceName, false)
	}

}
