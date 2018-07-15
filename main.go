package main

import (
	"github.com/modulus/goRabbitReceiver/receive"
)

func main() {

	receiver := receive.NewReceiver("config.json")
	receiver.ReadMessages()

}
