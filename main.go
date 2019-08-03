/*
$ go install
$ ~/bin/remind-me
*/

package main

import (
	"github.com/davidfloyd91/remind-me/db"
	"github.com/davidfloyd91/remind-me/server"
)

func main() {
	db.Connect()
	server.Start()
}
