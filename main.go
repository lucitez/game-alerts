package main

import (
	"fmt"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/lucitez/game-alerts/gamealerts"
	"github.com/lucitez/game-alerts/internal/logger"
)

func init() {
	functions.CloudEvent("SendGameAlert", gamealerts.SendGameAlert)

	logger.Init()
}

func main() {
	fmt.Println("hello")
}
