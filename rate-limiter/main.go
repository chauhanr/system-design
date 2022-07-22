package main

import (
	"net/http"

	"github.com/chauhanr/system-design/rate-limiter/api/app"
)

func main() {
	client := http.Client{}

	app.Startup(&client)
}
